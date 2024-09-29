package test

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
)

func GenerateRandomUser() *models.User {
	return &models.User{
		ID:    uuid.New(),
		Email: "john.doe@example.com",
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
