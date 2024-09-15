package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
)

func CreateSession(rdb *redis.Client, token *auth.AuthToken, userID uuid.UUID) error {
	s := models.UserSession{UserID: userID}
	sessionData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err := rdb.Set(
		context.Background(),
		token.TokenValue,
		sessionData,
		token.ValidDuration,
	).Err(); err != nil {
		return err
	}

	return nil
}

func GetSession(rdb *redis.Client, tokenValue string) (*models.UserSession, error) {
	sessionString, err := rdb.Get(context.Background(), tokenValue).Result()
	if err != nil {
		return nil, err
	}

	var s models.UserSession
	if err := json.Unmarshal([]byte(sessionString), &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func DeleteSession(rdb *redis.Client, tokenValue string) error {
	if err := rdb.Del(context.Background(), tokenValue).Err(); err != nil {
		return err
	}

	return nil
}

func DeleteUserSessions(rdb *redis.Client, userID uuid.UUID) error {
	var cursor uint64
	var keys []string
	var err error

	ctx := context.Background()

	for {
		keys, cursor, err = rdb.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			return fmt.Errorf("error scanning keys: %v", err)
		}

		for _, key := range keys {
			sessionData, err := rdb.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
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
				err := rdb.Del(ctx, key).Err()
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
