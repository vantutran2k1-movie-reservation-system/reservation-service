package auth

import "golang.org/x/crypto/bcrypt"

func GenerateHashedPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashAndPassword(hashedPassword, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}
