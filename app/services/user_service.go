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
	err = transaction.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		if err := s.userRepo.CreateUser(tx, &u); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
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

	token, err := auth.GenerateJwtToken(u.ID)
	if err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	if _, err := s.loginTokenRepo.GetActiveLoginToken(token.TokenValue); err != nil {
		if errors.IsRecordNotFoundError(err) {
			return "", errors.InternalServerError("Token %s already exists", token.TokenValue)
		}

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

// var LogoutUser = func(db *gorm.DB, rdb *redis.Client, tokenValue string) *errors.ApiError {
// 	if err := DeleteSession(rdb, tokenValue); err != nil {
// 		return errors.InternalServerError(err.Error())
// 	}

// 	if err := RevokeLoginToken(db, tokenValue); err != nil {
// 		return err
// 	}

// 	return nil
// }

// var UpdatePassword = func(db *gorm.DB, rdb *redis.Client, userID uuid.UUID, newPassword string) *errors.ApiError {
// 	var u models.User
// 	if err := db.Where(&models.User{ID: userID}).First(&u).Error; err != nil {
// 		return errors.InternalServerError(err.Error())
// 	}

// 	if err := auth.CompareHashAndPassword(u.PasswordHash, newPassword); err == nil {
// 		return errors.BadRequestError("New password can not be the same as current value")
// 	}

// 	p, err := auth.GenerateHashedPassword(newPassword)
// 	if err != nil {
// 		return errors.InternalServerError(err.Error())
// 	}

// 	u.PasswordHash = p
// 	u.UpdatedAt = time.Now().UTC()

// 	if err := db.Save(&u).Error; err != nil {
// 		return errors.InternalServerError(err.Error())
// 	}

// 	if err := DeleteUserSessions(rdb, userID); err != nil {
// 		return errors.InternalServerError(err.Error())
// 	}

// 	if err := RevokeUserLoginTokens(db, userID); err != nil {
// 		return err
// 	}

// 	return nil
// }
