package utils

import (
	"fmt"
	"mime/multipart"
	"time"

	"math/rand"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomUser() *models.User {
	return &models.User{
		ID:           uuid.New(),
		Email:        generateRandomEmail(),
		PasswordHash: generateRandomHashedPassword(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
}

func GenerateRandomUserWithID(id uuid.UUID) *models.User {
	return &models.User{
		ID:           id,
		Email:        generateRandomEmail(),
		PasswordHash: generateRandomHashedPassword(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
}

func GenerateRandomCreateUserRequest() *payloads.CreateUserRequest {
	return &payloads.CreateUserRequest{
		Email:    generateRandomEmail(),
		Password: generateRandomPassword(),
	}
}

func GenerateRandomUpdatePasswordRequest() *payloads.UpdatePasswordRequest {
	return &payloads.UpdatePasswordRequest{
		Password: generateRandomPassword(),
	}
}

func GenerateRandomUserProfile() *models.UserProfile {
	phoneNumber := generateRandomPhoneNumber()
	dateOfBirth := generateRandomDate()
	profilePictureUrl := generateRandomURL()
	bio := generateRandomString(allChars, 50)

	return &models.UserProfile{
		ID:                uuid.New(),
		UserID:            uuid.New(),
		FirstName:         generateRandomName(),
		LastName:          generateRandomName(),
		PhoneNumber:       &phoneNumber,
		DateOfBirth:       &dateOfBirth,
		ProfilePictureUrl: &profilePictureUrl,
		Bio:               &bio,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
}

func GenerateRandomUserProfileWithIDAndUserID(id uuid.UUID, userID uuid.UUID) *models.UserProfile {
	phoneNumber := generateRandomPhoneNumber()
	dateOfBirth := generateRandomDate()
	profilePictureUrl := generateRandomURL()
	bio := generateRandomString(allChars, 50)

	return &models.UserProfile{
		ID:                id,
		UserID:            userID,
		FirstName:         generateRandomName(),
		LastName:          generateRandomName(),
		PhoneNumber:       &phoneNumber,
		DateOfBirth:       &dateOfBirth,
		ProfilePictureUrl: &profilePictureUrl,
		Bio:               &bio,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
}

func GenerateRandomCreateUserProfileRequest() *payloads.CreateUserProfileRequest {
	phoneNumber := generateRandomPhoneNumber()
	dateOfBirth := generateRandomDate()

	return &payloads.CreateUserProfileRequest{
		FirstName:   generateRandomName(),
		LastName:    generateRandomName(),
		PhoneNumber: &phoneNumber,
		DateOfBirth: &dateOfBirth,
	}
}

func GenerateRandomUpdateUserProfileRequest() *payloads.UpdateUserProfileRequest {
	phoneNumber := generateRandomPhoneNumber()
	dateOfBirth := generateRandomDate()

	return &payloads.UpdateUserProfileRequest{
		FirstName:   generateRandomName(),
		LastName:    generateRandomName(),
		PhoneNumber: &phoneNumber,
		DateOfBirth: &dateOfBirth,
	}
}

func GenerateRandomAuthToken() *auth.AuthToken {
	return &auth.AuthToken{
		TokenValue:    uuid.NewString(),
		CreatedAt:     time.Now(),
		ValidDuration: time.Duration(60 * time.Minute),
	}
}

func GenerateRandomLoginToken() *models.LoginToken {
	return &models.LoginToken{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		TokenValue: uuid.NewString(),
		CreatedAt:  time.Now().UTC(),
		ExpiresAt:  time.Now().UTC().Add(60 * time.Minute),
	}
}

func GenerateRandomUserSession() *models.UserSession {
	return &models.UserSession{
		UserID: uuid.New(),
		Email:  generateRandomEmail(),
	}
}

func GenerateRandomFileHeader() *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: generateRandomName(),
		Size:     100,
		Header:   map[string][]string{constants.CONTENT_TYPE: {constants.IMAGE_PNG}},
	}
}

func GenerateRandomMovie() *models.Movie {
	description := generateRandomString(letterChars, 100)
	language := generateRandomString(letterChars, 10)
	rating := generateRandomFloat(0, 5)

	return &models.Movie{
		ID:              uuid.New(),
		Title:           generateRandomString(letterChars, 10),
		Description:     &description,
		ReleaseDate:     generateRandomDate(),
		DurationMinutes: generateRandomInt(100, 200),
		Language:        &language,
		Rating:          &rating,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		CreatedBy:       uuid.New(),
		LastUpdatedBy:   uuid.New(),
	}
}

func GenerateRandomResponseMeta() *models.ResponseMeta {
	prevUrl := generateRandomURL()
	nextUrl := generateRandomURL()

	return &models.ResponseMeta{
		Limit:   generateRandomInt(1, 10),
		Offset:  generateRandomInt(0, 5),
		Total:   generateRandomInt(20, 30),
		NextUrl: &nextUrl,
		PrevUrl: &prevUrl,
	}
}

func GenerateRandomCreateMovieRequest() *payloads.CreateMovieRequest {
	description := generateRandomString(allChars, 100)
	language := generateRandomString(lowercaseChars, 10)
	rating := generateRandomFloat(0, 5)

	return &payloads.CreateMovieRequest{
		Title:           generateRandomString(allChars, 10),
		Description:     &description,
		ReleaseDate:     generateRandomDate(),
		DurationMinutes: generateRandomInt(100, 200),
		Language:        &language,
		Rating:          &rating,
	}
}

const lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
const numberChars = "0123456789"
const letterChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

var urlProtocols = []string{"http", "https"}
var urlDomains = []string{"example.com", "testsite.com", "mywebsite.org", "randomsite.net"}
var emailDomains = []string{"gmail.com", "yahoo.com", "outlook.com", "example.com"}

func generateRandomString(chars string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func generateRandomName() string {
	return generateRandomString(lowercaseChars, 5)
}

func generateRandomEmail() string {
	nameLength := rand.Intn(6) + 5
	domain := emailDomains[rand.Intn(len(emailDomains))]

	return fmt.Sprintf("%s@%s", generateRandomString(lowercaseChars, nameLength), domain)
}

func generateRandomPassword() string {
	return generateRandomString(allChars, 12)
}

func generateRandomHashedPassword() string {
	password := generateRandomPassword()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword)
}

func generateRandomPhoneNumber() string {
	return generateRandomString(numberChars, 10)
}

func generateRandomDate() string {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now().UTC()

	randomDays := rand.Int63n(int64(end.Sub(start).Hours() / 24))

	randomDate := start.AddDate(0, 0, int(randomDays)-1)

	return randomDate.Format("2006-01-02")
}

func generateRandomURL() string {
	protocol := urlProtocols[rand.Intn(len(urlProtocols))]
	subdomain := generateRandomString(lowercaseChars, rand.Intn(5)+3)
	domain := urlDomains[rand.Intn(len(urlDomains))]
	path := generateRandomString(lowercaseChars+"/", rand.Intn(10)+5)
	queryParams := generateRandomString(lowercaseChars, rand.Intn(3)+3) + "=value"

	return fmt.Sprintf("%s://%s.%s/%s?%s", protocol, subdomain, domain, path, queryParams)
}

func generateRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func generateRandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
