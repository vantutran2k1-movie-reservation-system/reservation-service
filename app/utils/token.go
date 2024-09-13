package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthToken struct {
	TokenValue    string
	CreatedAt     time.Time
	ValidDuration time.Duration
}

type JwtClaims struct {
	UserID uuid.UUID
	jwt.RegisteredClaims
}

func GenerateJwtToken(userID uuid.UUID) (*AuthToken, error) {
	jwtExpiresAfterStr := os.Getenv("JWT_TOKEN_EXPIRES_AFTER_MINUTES")
	jwtExpiresAfter, err := strconv.Atoi(jwtExpiresAfterStr)
	if err != nil {
		return nil, fmt.Errorf("invalid token expiry minutes: %v", err)
	}

	validDuration := time.Duration(jwtExpiresAfter) * time.Minute

	claims := JwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(validDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	t := AuthToken{
		TokenValue:    tokenString,
		CreatedAt:     time.Now().UTC(),
		ValidDuration: validDuration,
	}

	return &t, nil
}
