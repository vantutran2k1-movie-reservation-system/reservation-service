package services

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(id uuid.UUID, includeProfile bool) (*models.User, *errors.ApiError)
	UserExistsByEmail(email string) (bool, *errors.ApiError)
	CreateUser(req payloads.CreateUserRequest) (*models.User, *errors.ApiError)
	LoginUser(req payloads.LoginUserRequest) (*models.LoginToken, *errors.ApiError)
	LogoutUser(tokenValue string) *errors.ApiError
	UpdateUserPassword(userID uuid.UUID, req payloads.UpdatePasswordRequest) *errors.ApiError
	CreatePasswordResetToken(req payloads.CreatePasswordResetTokenRequest) (*models.PasswordResetToken, *errors.ApiError)
	ResetUserPassword(resetToken string, request payloads.ResetPasswordRequest) *errors.ApiError
}

type userService struct {
	db                        *gorm.DB
	rdb                       *redis.Client
	authenticator             auth.Authenticator
	transactionManager        transaction.TransactionManager
	userRepo                  repositories.UserRepository
	userProfileRepo           repositories.UserProfileRepository
	loginTokenRepo            repositories.LoginTokenRepository
	userSessionRepo           repositories.UserSessionRepository
	passwordResetTokenRepo    repositories.PasswordResetTokenRepository
	userRegistrationTokenRepo repositories.UserRegistrationTokenRepository
	notificationRepo          repositories.NotificationRepository
}

func NewUserService(
	db *gorm.DB,
	rdb *redis.Client,
	authenticator auth.Authenticator,
	transactionManager transaction.TransactionManager,
	userRepo repositories.UserRepository,
	userProfileRepo repositories.UserProfileRepository,
	loginTokenRepo repositories.LoginTokenRepository,
	userSessionRepo repositories.UserSessionRepository,
	passwordResetTokenRepo repositories.PasswordResetTokenRepository,
	userRegistrationTokenRepo repositories.UserRegistrationTokenRepository,
	notificationRepo repositories.NotificationRepository,
) UserService {
	return &userService{
		db:                        db,
		rdb:                       rdb,
		authenticator:             authenticator,
		transactionManager:        transactionManager,
		userRepo:                  userRepo,
		userProfileRepo:           userProfileRepo,
		loginTokenRepo:            loginTokenRepo,
		userSessionRepo:           userSessionRepo,
		passwordResetTokenRepo:    passwordResetTokenRepo,
		userRegistrationTokenRepo: userRegistrationTokenRepo,
		notificationRepo:          notificationRepo,
	}
}

func (s *userService) GetUser(id uuid.UUID, includeProfile bool) (*models.User, *errors.ApiError) {
	u, err := s.getUserById(id, includeProfile)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if u == nil {
		return nil, errors.NotFoundError("user does not exist")
	}

	return u, nil
}

func (s *userService) UserExistsByEmail(email string) (bool, *errors.ApiError) {
	exist, err := s.userRepo.UserExists(filters.UserFilter{
		Filter: &filters.SingleFilter{},
		Email:  &filters.Condition{Operator: filters.OpEqual, Value: email},
	})
	if err != nil {
		return false, errors.InternalServerError(err.Error())
	}
	return exist, nil
}

func (s *userService) CreateUser(req payloads.CreateUserRequest) (*models.User, *errors.ApiError) {
	u, err := s.getUserByEmail(req.Email, false)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if u != nil {
		return nil, errors.BadRequestError("email %s already exists", req.Email)
	}

	hashedPassword, err := s.authenticator.GenerateHashedPassword(req.Password)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	userId := uuid.New()
	currentTime := time.Now().UTC()
	u = &models.User{
		ID:           userId,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     false,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	p := &models.UserProfile{
		ID:          uuid.New(),
		UserID:      userId,
		FirstName:   req.Profile.FirstName,
		LastName:    req.Profile.LastName,
		PhoneNumber: req.Profile.PhoneNumber,
		DateOfBirth: req.Profile.DateOfBirth,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}
	t := &models.UserRegistrationToken{
		ID:         uuid.New(),
		UserID:     userId,
		TokenValue: s.authenticator.GenerateRegistrationToken(),
		IsUsed:     false,
		CreatedAt:  currentTime,
		ExpiresAt:  currentTime,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		if err := s.userRepo.CreateUser(tx, u); err != nil {
			return err
		}
		if err := s.userProfileRepo.CreateUserProfile(tx, p); err != nil {
			return err
		}

		return s.userRegistrationTokenRepo.CreateToken(tx, t)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	e := payloads.UserRegistrationEvent{
		Email:             req.Email,
		FirstName:         req.Profile.FirstName,
		LastName:          req.Profile.LastName,
		VerificationToken: t.TokenValue,
		CreatedAt:         currentTime,
	}
	if err := s.notificationRepo.SendUserRegistrationEvent(e); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return u, nil
}

func (s *userService) LoginUser(req payloads.LoginUserRequest) (*models.LoginToken, *errors.ApiError) {
	u, err := s.getUserByEmail(req.Email, false)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if u == nil {
		return nil, errors.UnauthorizedError("invalid email %s", req.Email)
	}

	if !s.authenticator.DoPasswordsMatch(u.PasswordHash, req.Password) {
		return nil, errors.UnauthorizedError("invalid password")
	}

	token := s.authenticator.GenerateLoginToken()

	t, err := s.loginTokenRepo.GetLoginToken(filters.LoginTokenFilter{
		Filter:     &filters.SingleFilter{Logic: filters.And},
		TokenValue: &filters.Condition{Operator: filters.OpEqual, Value: token},
	})
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if t != nil {
		return nil, errors.InternalServerError("token value already exists")
	}

	now := time.Now().UTC()
	tokenExpiresAfter, err := strconv.Atoi(os.Getenv("LOGIN_TOKEN_EXPIRES_AFTER_MINUTES"))
	if err != nil {
		return nil, errors.InternalServerError("invalid token expiry time: %v", err)
	}
	validDuration := time.Duration(tokenExpiresAfter) * time.Minute

	t = &models.LoginToken{
		ID:         uuid.New(),
		UserID:     u.ID,
		TokenValue: token,
		CreatedAt:  now,
		ExpiresAt:  now.Add(validDuration),
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.loginTokenRepo.CreateLoginToken(tx, t)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInRedisTransaction(s.rdb, func(tx *redis.Tx) error {
		return s.userSessionRepo.CreateUserSession(
			s.userSessionRepo.GetUserSessionID(token),
			validDuration,
			&models.UserSession{UserID: u.ID, Email: u.Email},
		)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return t, nil
}

func (s *userService) LogoutUser(tokenValue string) *errors.ApiError {
	token, err := s.loginTokenRepo.GetLoginToken(filters.LoginTokenFilter{
		Filter:     &filters.SingleFilter{Logic: filters.And},
		TokenValue: &filters.Condition{Operator: filters.OpEqual, Value: tokenValue},
	})
	if err != nil {
		return errors.InternalServerError(err.Error())

	}
	if token == nil {
		return errors.BadRequestError("token not found")
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.loginTokenRepo.RevokeLoginToken(tx, token)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInRedisTransaction(s.rdb, func(tx *redis.Tx) error {
		return s.userSessionRepo.DeleteUserSession(s.userSessionRepo.GetUserSessionID(token.TokenValue))
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (s *userService) UpdateUserPassword(userID uuid.UUID, req payloads.UpdatePasswordRequest) *errors.ApiError {
	u, err := s.getUserById(userID, false)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	if u == nil {
		return errors.NotFoundError("user does not exist")
	}

	if s.authenticator.DoPasswordsMatch(u.PasswordHash, req.Password) {
		return errors.BadRequestError("new password can not be the same as current value")
	}

	p, err := s.authenticator.GenerateHashedPassword(req.Password)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.updateUserPassword(tx, u, p)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInRedisTransaction(s.rdb, func(tx *redis.Tx) error {
		return s.userSessionRepo.DeleteUserSessions(userID)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (s *userService) CreatePasswordResetToken(req payloads.CreatePasswordResetTokenRequest) (*models.PasswordResetToken, *errors.ApiError) {
	u, err := s.getUserByEmail(req.Email, false)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if u == nil {
		return nil, errors.NotFoundError("email %s does not exist", req.Email)
	}

	token := s.authenticator.GeneratePasswordResetToken()

	t, err := s.passwordResetTokenRepo.GetToken(filters.PasswordResetTokenFilter{
		Filter:     &filters.SingleFilter{},
		TokenValue: &filters.Condition{Operator: filters.OpEqual, Value: token},
	})
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if t != nil {
		return nil, errors.InternalServerError("token value already exists")
	}

	tokenExpiresAfter, err := strconv.Atoi(os.Getenv("PASSWORD_RESET_TOKEN_EXPIRES_AFTER_MINUTES"))
	if err != nil {
		return nil, errors.InternalServerError("invalid token expiry time: %v", err)
	}

	now := time.Now().UTC()
	t = &models.PasswordResetToken{
		ID:         uuid.New(),
		UserID:     u.ID,
		TokenValue: token,
		IsUsed:     false,
		CreatedAt:  now,
		ExpiresAt:  now.Add(time.Duration(tokenExpiresAfter) * time.Minute),
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.passwordResetTokenRepo.CreateToken(tx, t)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return t, nil
}

func (s *userService) ResetUserPassword(resetToken string, req payloads.ResetPasswordRequest) *errors.ApiError {
	t, err := s.passwordResetTokenRepo.GetToken(filters.PasswordResetTokenFilter{
		Filter:     &filters.SingleFilter{},
		TokenValue: &filters.Condition{Operator: filters.OpEqual, Value: resetToken},
	})
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	if t == nil {
		return errors.UnauthorizedError("invalid or expired token")
	}

	u, err := s.getUserById(t.UserID, false)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	if u == nil {
		return errors.InternalServerError("user not found")
	}

	p, err := s.authenticator.GenerateHashedPassword(req.Password)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		if err := s.updateUserPassword(tx, u, p); err != nil {
			return err
		}

		if err := s.passwordResetTokenRepo.UseToken(tx, t); err != nil {
			return err
		}

		tokens, err := s.getRemainingUserActivePasswordResetTokens(t.UserID, t.TokenValue)
		if err != nil {
			return err
		}
		if err := s.passwordResetTokenRepo.RevokeTokens(tx, tokens); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInRedisTransaction(s.rdb, func(tx *redis.Tx) error {
		return s.userSessionRepo.DeleteUserSessions(u.ID)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (s *userService) getUserById(id uuid.UUID, includeProfile bool) (*models.User, error) {
	return s.userRepo.GetUser(filters.UserFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: id},
	}, includeProfile)
}

func (s *userService) getUserByEmail(email string, includeProfile bool) (*models.User, error) {
	return s.userRepo.GetUser(filters.UserFilter{
		Filter: &filters.SingleFilter{},
		Email:  &filters.Condition{Operator: filters.OpEqual, Value: email},
	}, includeProfile)
}

func (s *userService) getRemainingUserActivePasswordResetTokens(userID uuid.UUID, tokenValue string) ([]*models.PasswordResetToken, error) {
	allTokens, err := s.passwordResetTokenRepo.GetTokens(filters.PasswordResetTokenFilter{
		Filter:    &filters.MultiFilter{},
		UserID:    &filters.Condition{Operator: filters.OpEqual, Value: userID},
		ExpiresAt: &filters.Condition{Operator: filters.OpGreater, Value: time.Now().UTC()},
	})
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	var tokens []*models.PasswordResetToken
	for _, token := range allTokens {
		if token.TokenValue != tokenValue {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (s *userService) updateUserPassword(tx *gorm.DB, u *models.User, p string) error {
	user, err := s.userRepo.UpdatePassword(tx, u, p)
	if err != nil {
		return err
	}

	u = user

	if err := s.loginTokenRepo.RevokeUserLoginTokens(tx, u.ID); err != nil {
		return err
	}

	return nil
}
