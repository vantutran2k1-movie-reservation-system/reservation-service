package sessions

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

type Session struct {
	UserID uuid.UUID
}

func CreateSession(token *utils.AuthToken, userID uuid.UUID) error {
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
