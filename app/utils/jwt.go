package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(userID uuid.UUID) (string, error) {
	jwtExpiresAfterStr := os.Getenv("JWT_TOKEN_EXPIRES_AFTER_MINUTES")
	jwtExpiresAfter, err := strconv.Atoi(jwtExpiresAfterStr)
	if err != nil {
		return "", fmt.Errorf("invalid token expiry minutes: %v", err)
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(jwtExpiresAfter) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
