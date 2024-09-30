package auth

import "golang.org/x/crypto/bcrypt"

type Authenticator interface {
	GenerateHashedPassword(rawPassword string) (string, error)
	IsPasswordsMatch(hashedPassword, rawPassword string) bool
}

func NewAuthenticator() Authenticator {
	return &bcryptAuthenticator{}
}

type bcryptAuthenticator struct{}

func (a *bcryptAuthenticator) GenerateHashedPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (a *bcryptAuthenticator) IsPasswordsMatch(hashedPassword, rawPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword)); err != nil {
		return false
	}

	return true
}
