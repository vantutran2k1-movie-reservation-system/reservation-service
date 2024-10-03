package utils

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
)

func GenerateSampleUser() *models.User {
	return &models.User{
		ID:           uuid.New(),
		Email:        "email@example.com",
		PasswordHash: "Hashed password",
	}
}

func GenerateSampleUserProfile() *models.UserProfile {
	phoneNumber := "1234567890"
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
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
}

func GenerateSampleLoginToken() *models.LoginToken {
	return &models.LoginToken{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		TokenValue: "sample_token_value",
		CreatedAt:  time.Now().UTC(),
		ExpiresAt:  time.Now().UTC().Add(60 * time.Minute),
	}
}

func GenerateSampleUserSession() *models.UserSession {
	return &models.UserSession{
		UserID: uuid.New(),
		Email:  "email@example.com",
	}
}

func GenerateSampleAuthToken() *auth.AuthToken {
	return &auth.AuthToken{
		TokenValue:    "sample_token_value",
		CreatedAt:     time.Now(),
		ValidDuration: time.Duration(60 * time.Minute),
	}
}

func GenerateSampleFileHeader() *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: "test-image.png",
		Size:     12345,
		Header:   map[string][]string{"Content-Type": {"image/png"}},
	}
}

func GenerateSampleCreateUserRequest() *payloads.CreateUserRequest {
	return &payloads.CreateUserRequest{
		Email:    "email@example.com",
		Password: "password",
	}
}

func GenerateSampleUpdatePasswordRequest() *payloads.UpdatePasswordRequest {
	return &payloads.UpdatePasswordRequest{
		Password: "new password",
	}
}

func GenerateSampleCreateUserProfileRequest() *payloads.CreateUserProfileRequest {
	phoneNumber := "1234567890"
	dateOfBirth := "1990-01-01"

	return &payloads.CreateUserProfileRequest{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: &phoneNumber,
		DateOfBirth: &dateOfBirth,
	}
}

func GenerateSampleUpdateUserProfileRequest() *payloads.UpdateUserProfileRequest {
	phoneNumber := "1357924680"
	dateOfBirth := "2000-01-01"

	return &payloads.UpdateUserProfileRequest{
		FirstName:   "First",
		LastName:    "Last",
		PhoneNumber: &phoneNumber,
		DateOfBirth: &dateOfBirth,
	}
}

func GenerateSampleMovie() *models.Movie {
	description := "Movie description"
	language := "English"
	rating := 4.5

	return &models.Movie{
		ID:              uuid.New(),
		Title:           "Movie title",
		Description:     &description,
		ReleaseDate:     "2024-01-01",
		DurationMinutes: 120,
		Language:        &language,
		Rating:          &rating,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		CreatedBy:       uuid.New(),
		LastUpdatedBy:   uuid.New(),
	}
}

func GenerateSampleCreateMovieRequest() *payloads.CreateMovieRequest {
	description := "A thrilling story of a group of explorers embarking on a dangerous mission to uncover a lost city in the Amazon rainforest."
	language := "English"
	rating := 3.0

	return &payloads.CreateMovieRequest{
		Title:           "The Last Adventure",
		Description:     &description,
		ReleaseDate:     "2024-01-01",
		DurationMinutes: 120,
		Language:        &language,
		Rating:          &rating,
	}
}
