package services

import (
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
	GetUser(userID uuid.UUID) (*models.User, *errors.ApiError)
	CreateUser(email string, password string) (*models.User, *errors.ApiError)
	LoginUser(email string, password string) (*models.LoginToken, *errors.ApiError)
	LogoutUser(tokenValue string) *errors.ApiError
	UpdateUserPassword(userID uuid.UUID, password string) *errors.ApiError
	CreatePasswordResetToken(email string) (*models.PasswordResetToken, *errors.ApiError)
}

type userService struct {
	db                     *gorm.DB
	rdb                    *redis.Client
	authenticator          auth.Authenticator
	transactionManager     transaction.TransactionManager
	userRepo               repositories.UserRepository
	loginTokenRepo         repositories.LoginTokenRepository
	userSessionRepo        repositories.UserSessionRepository
	passwordResetTokenRepo repositories.PasswordResetTokenRepository
}

func NewUserService(
	db *gorm.DB,
	rdb *redis.Client,
	authenticator auth.Authenticator,
	transactionManager transaction.TransactionManager,
	userRepo repositories.UserRepository,
	loginTokenRepo repositories.LoginTokenRepository,
	userSessionRepo repositories.UserSessionRepository,
	passwordResetTokenRepo repositories.PasswordResetTokenRepository,
) UserService {
	return &userService{
		db:                     db,
		rdb:                    rdb,
		authenticator:          authenticator,
		transactionManager:     transactionManager,
		userRepo:               userRepo,
		loginTokenRepo:         loginTokenRepo,
		userSessionRepo:        userSessionRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
	}
}

func (s *userService) GetUser(userID uuid.UUID) (*models.User, *errors.ApiError) {
	u, err := s.userRepo.GetUser(userID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, errors.BadRequestError("User does not exist")
		}

		return nil, errors.InternalServerError(err.Error())
	}

	return u, nil
}

func (s *userService) CreateUser(email string, password string) (*models.User, *errors.ApiError) {
	_, err := s.userRepo.FindUserByEmail(email)

	if err == nil {
		return nil, errors.BadRequestError("Email %s already exists", email)
	}

	if !errors.IsRecordNotFoundError(err) {
		return nil, errors.InternalServerError(err.Error())
	}

	hashedPassword, err := s.authenticator.GenerateHashedPassword(password)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	u := models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.userRepo.CreateUser(tx, &u)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &u, nil
}

func (s *userService) LoginUser(email string, password string) (*models.LoginToken, *errors.ApiError) {
	u, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, errors.UnauthorizedError("Invalid email %s", email)
		}

		return nil, errors.InternalServerError(err.Error())
	}

	if !s.authenticator.DoPasswordsMatch(u.PasswordHash, password) {
		return nil, errors.UnauthorizedError("Invalid password")
	}

	token := s.authenticator.GenerateLoginToken()
	_, err = s.loginTokenRepo.GetActiveLoginToken(token)
	if err == nil {
		return nil, errors.InternalServerError("Token value already exists")
	}
	if !errors.IsRecordNotFoundError(err) {
		return nil, errors.InternalServerError(err.Error())
	}

	now := time.Now().UTC()
	tokenExpiresAfter, err := strconv.Atoi(os.Getenv("AUTH_TOKEN_EXPIRES_AFTER_MINUTES"))
	if err != nil {
		return nil, errors.InternalServerError("invalid token expiry minutes: %v", err)
	}
	validDuration := time.Duration(tokenExpiresAfter) * time.Minute

	t := models.LoginToken{
		ID:         uuid.New(),
		UserID:     u.ID,
		TokenValue: token,
		CreatedAt:  now,
		ExpiresAt:  now.Add(validDuration),
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.loginTokenRepo.CreateLoginToken(tx, &t)
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

	return &t, nil
}

func (s *userService) LogoutUser(tokenValue string) *errors.ApiError {
	token, err := s.loginTokenRepo.GetActiveLoginToken(tokenValue)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return errors.BadRequestError("Token not found")
		}

		return errors.InternalServerError(err.Error())
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

func (s *userService) UpdateUserPassword(userID uuid.UUID, password string) *errors.ApiError {
	u, err := s.userRepo.GetUser(userID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return errors.BadRequestError("User does not exist")
		}

		return errors.InternalServerError(err.Error())
	}

	if s.authenticator.DoPasswordsMatch(u.PasswordHash, password) {
		return errors.BadRequestError("New password can not be the same as current value")
	}

	p, err := s.authenticator.GenerateHashedPassword(password)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		user, err := s.userRepo.UpdatePassword(tx, u, p)
		if err != nil {
			return err
		}

		u = user

		if err := s.loginTokenRepo.RevokeUserLoginTokens(tx, u.ID); err != nil {
			return err
		}

		return nil
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

func (s *userService) CreatePasswordResetToken(email string) (*models.PasswordResetToken, *errors.ApiError) {
	u, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, errors.NotFoundError("email does not exist")
		}

		return nil, errors.InternalServerError(err.Error())
	}

	token := s.authenticator.GeneratePasswordResetToken()
	_, err = s.passwordResetTokenRepo.GetActivePasswordResetToken(token)
	if err == nil {
		return nil, errors.InternalServerError("token value already exists")
	}
	if !errors.IsRecordNotFoundError(err) {
		return nil, errors.InternalServerError(err.Error())
	}

	tokenExpiresAfter, err := strconv.Atoi(os.Getenv("PASSWORD_RESET_TOKEN_EXPIRES_AFTER_MINUTES"))
	if err != nil {
		return nil, errors.InternalServerError("invalid token expiry minutes: %v", err)
	}

	now := time.Now().UTC()
	t := models.PasswordResetToken{
		ID:         uuid.New(),
		UserID:     u.ID,
		TokenValue: token,
		IsUsed:     false,
		CreatedAt:  now,
		ExpiresAt:  now.Add(time.Duration(tokenExpiresAfter) * time.Minute),
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.passwordResetTokenRepo.CreateToken(tx, &t)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &t, nil
}
