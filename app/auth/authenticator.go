package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	GenerateHashedPassword(rawPassword string) (string, error)
	DoPasswordsMatch(hashedPassword, rawPassword string) bool
	GenerateLoginToken() string
	GeneratePasswordResetToken() string
	GenerateRegistrationToken() string
}

func NewAuthenticator() Authenticator {
	return &authenticator{}
}

type authenticator struct{}

func (a *authenticator) GenerateHashedPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (a *authenticator) DoPasswordsMatch(hashedPassword, rawPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword)); err != nil {
		return false
	}

	return true
}

func (a *authenticator) GenerateLoginToken() string {
	return uuid.NewString()
}

func (a *authenticator) GeneratePasswordResetToken() string { return uuid.NewString() }

func (a *authenticator) GenerateRegistrationToken() string { return uuid.NewString() }
