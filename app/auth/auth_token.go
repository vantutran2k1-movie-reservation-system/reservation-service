package auth

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

type JwtTokenGenerator struct{}

type JwtClaims struct {
	jwt.RegisteredClaims
}

func (g *JwtTokenGenerator) GenerateToken() (*AuthToken, error) {
	jwtExpiresAfterStr := os.Getenv("AUTH_TOKEN_EXPIRES_AFTER_MINUTES")
	jwtExpiresAfter, err := strconv.Atoi(jwtExpiresAfterStr)
	if err != nil {
		return nil, fmt.Errorf("invalid token expiry minutes: %v", err)
	}

	validDuration := time.Duration(jwtExpiresAfter) * time.Minute

	claims := JwtClaims{
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
