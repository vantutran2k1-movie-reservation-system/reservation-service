package services

import (
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
	CreateUser(email string, password string) (*models.User, *errors.ApiError)
	LoginUser(email string, password string) (string, *errors.ApiError)
	LogoutUser(tokenValue string) *errors.ApiError
	UpdateUserPassword(userID uuid.UUID, password string) *errors.ApiError
}

type userService struct {
	db              *gorm.DB
	rdb             *redis.Client
	userRepo        repositories.UserRepository
	loginTokenRepo  repositories.LoginTokenRepository
	userSessionRepo repositories.UserSessionRepository
}

func NewUserService(
	db *gorm.DB,
	rdb *redis.Client,
	userRepo repositories.UserRepository,
	loginTokenRepo repositories.LoginTokenRepository,
	userSessionRepo repositories.UserSessionRepository,
) UserService {
	return &userService{
		db:              db,
		rdb:             rdb,
		userRepo:        userRepo,
		loginTokenRepo:  loginTokenRepo,
		userSessionRepo: userSessionRepo,
	}
}

func (s *userService) CreateUser(email string, password string) (*models.User, *errors.ApiError) {
	_, err := s.userRepo.FindUserByEmail(email)

	if err == nil {
		return nil, errors.BadRequestError("Email %s already exists", email)
	}

	if !errors.IsRecordNotFoundError(err) {
		return nil, errors.InternalServerError(err.Error())
	}

	hashedPassword, err := auth.GenerateHashedPassword(password)
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
	if err := transaction.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.userRepo.CreateUser(tx, &u)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &u, nil
}

func (s *userService) LoginUser(email string, password string) (string, *errors.ApiError) {
	u, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return "", errors.BadRequestError("Invalid email %s", email)
		}

		return "", errors.InternalServerError(err.Error())
	}

	if err := auth.CompareHashAndPassword(u.PasswordHash, password); err != nil {
		return "", errors.BadRequestError("Invalid password")
	}

	token, err := auth.GenerateBasicToken(u.ID)
	if err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	_, err = s.loginTokenRepo.GetActiveLoginToken(token.TokenValue)
	if err == nil {
		return "", errors.InternalServerError("Token value already exists")
	}
	if !errors.IsRecordNotFoundError(err) {
		return "", errors.InternalServerError(err.Error())
	}

	if err := transaction.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.loginTokenRepo.CreateLoginToken(tx, &models.LoginToken{
			ID:         uuid.New(),
			UserID:     u.ID,
			TokenValue: token.TokenValue,
			CreatedAt:  token.CreatedAt,
			ExpiresAt:  token.CreatedAt.Add(token.ValidDuration),
		})
	}); err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	if err := transaction.ExecuteInRedisTransaction(s.rdb, func(tx *redis.Tx) error {
		return s.userSessionRepo.CreateUserSession(
			s.userSessionRepo.GetUserSessionID(token.TokenValue),
			token.ValidDuration,
			&models.UserSession{UserID: u.ID},
		)
	}); err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	return token.TokenValue, nil
}

func (s *userService) LogoutUser(tokenValue string) *errors.ApiError {
	token, err := s.loginTokenRepo.GetActiveLoginToken(tokenValue)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			return errors.BadRequestError("Token not found")
		}

		return errors.InternalServerError(err.Error())
	}

	if err := transaction.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.loginTokenRepo.RevokeLoginToken(tx, token)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := transaction.ExecuteInRedisTransaction(s.rdb, func(tx *redis.Tx) error {
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
			return errors.BadRequestError("Invalid user id")
		}

		return errors.InternalServerError(err.Error())
	}

	if err := auth.CompareHashAndPassword(u.PasswordHash, password); err == nil {
		return errors.BadRequestError("New password can not be the same as current value")
	}

	p, err := auth.GenerateHashedPassword(password)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}

	u.PasswordHash = p
	u.UpdatedAt = time.Now().UTC()

	if err := transaction.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		if err := s.userRepo.UpdateUser(tx, u); err != nil {
			return err
		}

		if err := s.loginTokenRepo.RevokeUserLoginTokens(tx, u.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	if err := transaction.ExecuteInRedisTransaction(s.rdb, func(tx *redis.Tx) error {
		return s.userSessionRepo.DeleteUserSessions(userID)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}
