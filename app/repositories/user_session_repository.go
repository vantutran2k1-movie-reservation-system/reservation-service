package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
)

type UserSessionRepository interface {
	GetUserSession(sessionID string) (*models.UserSession, error)
	GetUserSessionID(tokenValue string) string
	CreateUserSession(sessionID string, expiration time.Duration, session *models.UserSession) error
	DeleteUserSession(sessionID string) error
	DeleteUserSessions(userID uuid.UUID) error
}

type userSessionRepository struct {
	ctx context.Context
	rdb *redis.Client
}

func NewUserSessionRepository(rdb *redis.Client) UserSessionRepository {
	return &userSessionRepository{ctx: context.Background(), rdb: rdb}
}

func (r *userSessionRepository) GetUserSession(sessionID string) (*models.UserSession, error) {
	sessionString, err := r.rdb.Get(r.ctx, sessionID).Result()
	if err != nil {
		return nil, err
	}

	var s models.UserSession
	if err := json.Unmarshal([]byte(sessionString), &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *userSessionRepository) GetUserSessionID(tokenValue string) string {
	return fmt.Sprintf("session:%s", tokenValue)
}

func (r *userSessionRepository) CreateUserSession(sessionID string, expiration time.Duration, session *models.UserSession) error {
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	if err := r.rdb.Set(r.ctx, sessionID, sessionData, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (r *userSessionRepository) DeleteUserSession(sessionID string) error {
	if err := r.rdb.Del(r.ctx, sessionID).Err(); err != nil {
		return err
	}

	return nil
}

func (r *userSessionRepository) DeleteUserSessions(userID uuid.UUID) error {
	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = r.rdb.Scan(r.ctx, cursor, "*", 100).Result()
		if err != nil {
			return fmt.Errorf("error scanning keys: %v", err)
		}

		for _, key := range keys {
			sessionData, err := r.rdb.Get(r.ctx, key).Result()
			if err != nil {
				if errors.IsRedisKeyNotFoundError(err) {
					continue
				}
				return fmt.Errorf("error getting key %s: %v", key, err)
			}

			var session models.UserSession
			err = json.Unmarshal([]byte(sessionData), &session)
			if err != nil {
				return fmt.Errorf("error unmarshaling session data: %v", err)
			}

			if session.UserID == userID {
				err := r.rdb.Del(r.ctx, key).Err()
				if err != nil {
					return fmt.Errorf("error deleting key %s: %v", key, err)
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
