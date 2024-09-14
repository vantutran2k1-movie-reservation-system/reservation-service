package sessions

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

type Session struct {
	UserID uuid.UUID `json:"user_id"`
}

func CreateSession(token *auth.AuthToken, userID uuid.UUID) error {
	s := Session{UserID: userID}
	sessionData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err := config.RedisClient.Set(
		context.Background(),
		token.TokenValue,
		sessionData,
		token.ValidDuration,
	).Err(); err != nil {
		return err
	}

	return nil
}

func GetSession(tokenValue string) (*Session, error) {
	sessionString, err := config.RedisClient.Get(context.Background(), tokenValue).Result()
	if err != nil {
		return nil, err
	}

	var s Session
	if err := json.Unmarshal([]byte(sessionString), &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func DeleteSession(tokenValue string) error {
	if err := config.RedisClient.Del(context.Background(), tokenValue).Err(); err != nil {
		return err
	}

	return nil
}
