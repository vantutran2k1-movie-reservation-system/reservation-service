package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

type Session struct {
	UserID uuid.UUID
}

func CreateSession(sessionID string, userID uuid.UUID) error {
	sessionExpiresAfterStr := os.Getenv("REDIS_SESSION_EXPIRES_AFTER_MINUTES")
	sessionExpiresAfter, err := strconv.Atoi(sessionExpiresAfterStr)
	if err != nil {
		return fmt.Errorf("invalid session expiry minutes: %v", err)
	}

	s := Session{UserID: userID}
	sessionData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err := config.RedisClient.Set(context.Background(), sessionID, sessionData, time.Duration(sessionExpiresAfter)*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}
