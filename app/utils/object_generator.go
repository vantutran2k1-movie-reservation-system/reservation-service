package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
)

func GenerateRandomUser() *models.User {
	return &models.User{
		ID:           uuid.New(),
		Email:        "john.doe@example.com",
		PasswordHash: "Hashed password",
	}
}

func GenerateRandomUserProfile() *models.UserProfile {
	phoneNumber := "123-456-7890"
	dateOfBirth := "1990-01-01"
	profilePictureUrl := "http://example.com/profile.jpg"
	bio := "This is a sample bio."

	return &models.UserProfile{
		ID:                uuid.New(),
		UserID:            uuid.New(),
		FirstName:         "John",
		LastName:          "Doe",
		PhoneNumber:       &phoneNumber,
		DateOfBirth:       &dateOfBirth,
		ProfilePictureUrl: &profilePictureUrl,
		Bio:               &bio,
	}
}

func GenerateRandomLoginToken() *models.LoginToken {
	return &models.LoginToken{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		TokenValue: "sample_token_value",
		CreatedAt:  time.Now().UTC(),
		ExpiresAt:  time.Now().UTC().Add(60 * time.Minute),
	}
}

func GenerateRandomUserSession() *models.UserSession {
	return &models.UserSession{
		UserID: uuid.New(),
		Email:  "john.doe@example.com",
	}
}

func GenerateRandomAuthToken() *auth.AuthToken {
	return &auth.AuthToken{
		TokenValue:    "sample_token_value",
		CreatedAt:     time.Now(),
		ValidDuration: time.Duration(60 * time.Minute),
	}
}
