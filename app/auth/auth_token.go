package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	TokenValue    string
	CreatedAt     time.Time
	ValidDuration time.Duration
}

type TokenGenerator interface {
	GenerateToken() (*AuthToken, error)
}

func NewTokenGenerator() TokenGenerator {
	return &uuidTokenGenerator{}
}

type uuidTokenGenerator struct{}

func (g *uuidTokenGenerator) GenerateToken() (*AuthToken, error) {
	tokenExpiresAfterStr := os.Getenv("AUTH_TOKEN_EXPIRES_AFTER_MINUTES")
	tokenExpiresAfter, err := strconv.Atoi(tokenExpiresAfterStr)
	if err != nil {
		return nil, fmt.Errorf("invalid token expiry minutes: %v", err)
	}

	validDuration := time.Duration(tokenExpiresAfter) * time.Minute
	t := AuthToken{
		TokenValue:    uuid.NewString(),
		CreatedAt:     time.Now().UTC(),
		ValidDuration: validDuration,
	}

	return &t, nil
}
