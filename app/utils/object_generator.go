package utils

import (
	"fmt"
	"mime/multipart"
	"time"

	"math/rand"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"golang.org/x/crypto/bcrypt"
)

// Common generated fields
func GenerateResponseMeta() *models.ResponseMeta {
	prevUrl := GenerateURL()
	nextUrl := GenerateURL()

	return &models.ResponseMeta{
		Limit:   generateInt(1, 10),
		Offset:  generateInt(0, 5),
		Total:   generateInt(20, 30),
		NextUrl: &nextUrl,
		PrevUrl: &prevUrl,
	}
}

func GenerateEmail() string {
	nameLength := rand.Intn(6) + 5
	domain := emailDomains[rand.Intn(len(emailDomains))]

	return fmt.Sprintf("%s@%s", generateString(lowercaseChars, nameLength), domain)
}

func GeneratePassword() string {
	return generateString(letterChars, 12)
}

func GenerateHashedPassword() string {
	password := GeneratePassword()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword)
}

func GeneratePhoneNumber() string {
	return generateString(numberChars, 10)
}

func GenerateURL() string {
	protocol := urlProtocols[rand.Intn(len(urlProtocols))]
	subdomain := generateString(lowercaseChars, rand.Intn(5)+3)
	domain := urlDomains[rand.Intn(len(urlDomains))]
	path := generateString(lowercaseChars+"/", rand.Intn(10)+5)
	queryParams := generateString(lowercaseChars, rand.Intn(3)+3) + "=value"

	return fmt.Sprintf("%s://%s.%s/%s?%s", protocol, subdomain, domain, path, queryParams)
}

func GenerateName() string {
	return generateString(lowercaseChars, 5)
}

func GenerateDate() string {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := generateCurrentTime()

	days := rand.Int63n(int64(end.Sub(start).Hours() / 24))

	date := start.AddDate(0, 0, int(days)-1)

	return date.Format("2006-01-02")
}

// Models and payloads
//
// User
func GenerateUser() *models.User {
	return &models.User{
		ID:           generateUUID(),
		Email:        GenerateEmail(),
		PasswordHash: GenerateHashedPassword(),
		CreatedAt:    generateCurrentTime(),
		UpdatedAt:    generateCurrentTime(),
	}
}

func GenerateCreateUserRequest() *payloads.CreateUserRequest {
	return &payloads.CreateUserRequest{
		Email:    GenerateEmail(),
		Password: GeneratePassword(),
	}
}

func GenerateUpdatePasswordRequest() *payloads.UpdatePasswordRequest {
	return &payloads.UpdatePasswordRequest{
		Password: GeneratePassword(),
	}
}

// User profile
func GenerateUserProfile() *models.UserProfile {
	phoneNumber := GeneratePhoneNumber()
	dateOfBirth := GenerateDate()
	profilePictureUrl := GenerateURL()
	bio := generateString(allChars, 50)

	return &models.UserProfile{
		ID:                generateUUID(),
		UserID:            generateUUID(),
		FirstName:         GenerateName(),
		LastName:          GenerateName(),
		PhoneNumber:       &phoneNumber,
		DateOfBirth:       &dateOfBirth,
		ProfilePictureUrl: &profilePictureUrl,
		Bio:               &bio,
		CreatedAt:         generateCurrentTime(),
		UpdatedAt:         generateCurrentTime(),
	}
}

func GenerateCreateUserProfileRequest() *payloads.CreateUserProfileRequest {
	phoneNumber := GeneratePhoneNumber()
	dateOfBirth := GenerateDate()

	return &payloads.CreateUserProfileRequest{
		FirstName:   GenerateName(),
		LastName:    GenerateName(),
		PhoneNumber: &phoneNumber,
		DateOfBirth: &dateOfBirth,
	}
}

func GenerateUpdateUserProfileRequest() *payloads.UpdateUserProfileRequest {
	phoneNumber := GeneratePhoneNumber()
	dateOfBirth := GenerateDate()

	return &payloads.UpdateUserProfileRequest{
		FirstName:   GenerateName(),
		LastName:    GenerateName(),
		PhoneNumber: &phoneNumber,
		DateOfBirth: &dateOfBirth,
	}
}

// Login token
func GenerateLoginToken() *models.LoginToken {
	return &models.LoginToken{
		ID:         generateUUID(),
		UserID:     generateUUID(),
		TokenValue: uuid.NewString(),
		CreatedAt:  generateCurrentTime(),
		ExpiresAt:  generateCurrentTime().Add(60 * time.Minute),
	}
}

// User session
func GenerateUserSession() *models.UserSession {
	return &models.UserSession{
		UserID: generateUUID(),
		Email:  GenerateEmail(),
	}
}

func GenerateSessionID() string {
	return uuid.NewString()
}

// File
func GenerateFileHeader() *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: GenerateName(),
		Size:     100,
		Header:   map[string][]string{constants.ContentType: {constants.ImagePng}},
	}
}

// Movie
func GenerateMovie() *models.Movie {
	description := generateString(letterChars, 100)
	language := generateString(letterChars, 10)
	rating := generateFloat(0, 5)

	return &models.Movie{
		ID:              generateUUID(),
		Title:           generateString(letterChars, 10),
		Description:     &description,
		ReleaseDate:     GenerateDate(),
		DurationMinutes: generateInt(100, 200),
		Language:        &language,
		Rating:          &rating,
		CreatedAt:       generateCurrentTime(),
		UpdatedAt:       generateCurrentTime(),
		CreatedBy:       generateUUID(),
		LastUpdatedBy:   generateUUID(),
	}
}

func GenerateMovies(count int) []*models.Movie {
	movies := make([]*models.Movie, count)
	for i := 0; i < count; i++ {
		movies[i] = GenerateMovie()
	}

	return movies
}

func GenerateCreateMovieRequest() *payloads.CreateMovieRequest {
	description := generateString(allChars, 100)
	language := generateString(lowercaseChars, 10)
	rating := generateFloat(0, 5)

	return &payloads.CreateMovieRequest{
		Title:           generateString(allChars, 10),
		Description:     &description,
		ReleaseDate:     GenerateDate(),
		DurationMinutes: generateInt(100, 200),
		Language:        &language,
		Rating:          &rating,
	}
}

func GenerateUpdateMovieGenresRequest() *payloads.UpdateMovieGenresRequest {
	return &payloads.UpdateMovieGenresRequest{
		GenreIDs: []uuid.UUID{generateUUID(), generateUUID()},
	}
}

// Genre
func GenerateGenre() *models.Genre {
	return &models.Genre{
		ID:   generateUUID(),
		Name: generateString(lowercaseChars, 10),
	}
}

func GenerateGenres(count int) []*models.Genre {
	genres := make([]*models.Genre, count)
	for i := 0; i < count; i++ {
		genres[i] = GenerateGenre()
	}

	return genres
}

func GenerateCreateGenreRequest() *payloads.CreateGenreRequest {
	return &payloads.CreateGenreRequest{
		Name: generateString(letterChars, 10),
	}
}

func GenerateUpdateGenreRequest() *payloads.UpdateGenreRequest {
	return &payloads.UpdateGenreRequest{
		Name: generateString(letterChars, 10),
	}
}

// Password reset token
func GeneratePasswordResetToken() *models.PasswordResetToken {
	return &models.PasswordResetToken{
		ID:         generateUUID(),
		UserID:     generateUUID(),
		TokenValue: uuid.NewString(),
		IsUsed:     generateBool(),
		CreatedAt:  generateCurrentTime(),
		ExpiresAt:  generateCurrentTime().Add(60 * time.Minute),
	}
}

func GeneratePasswordResetTokens(count int) []*models.PasswordResetToken {
	tokens := make([]*models.PasswordResetToken, count)
	for i := 0; i < count; i++ {
		tokens[i] = GeneratePasswordResetToken()
	}

	return tokens
}

func GenerateCreatePasswordResetTokenRequest() *payloads.CreatePasswordResetTokenRequest {
	return &payloads.CreatePasswordResetTokenRequest{
		Email: GenerateEmail(),
	}
}

// Country
func GenerateCountry() *models.Country {
	return &models.Country{
		ID:   generateUUID(),
		Name: generateString(lowercaseChars, 10),
		Code: generateString(uppercaseChars, 2),
	}
}

func GenerateCountries(count int) []*models.Country {
	countries := make([]*models.Country, count)
	for i := 0; i < count; i++ {
		countries[i] = GenerateCountry()
	}

	return countries
}

func GenerateCreateCountryRequest() *payloads.CreateCountryRequest {
	return &payloads.CreateCountryRequest{
		Name: generateString(lowercaseChars, 10),
		Code: generateString(uppercaseChars, 2),
	}
}

// State
func GenerateState() *models.State {
	code := generateString(uppercaseChars, 2)

	return &models.State{
		ID:        generateUUID(),
		Name:      generateString(lowercaseChars, 10),
		Code:      &code,
		CountryID: generateUUID(),
	}
}

func GenerateStates(count int) []*models.State {
	states := make([]*models.State, count)
	for i := 0; i < count; i++ {
		states[i] = GenerateState()
	}

	return states
}

func GenerateCreateStateRequest() *payloads.CreateStateRequest {
	code := generateString(uppercaseChars, 2)

	return &payloads.CreateStateRequest{
		Name: generateString(lowercaseChars, 10),
		Code: &code,
	}
}

// City
func GenerateCity() *models.City {
	return &models.City{
		ID:      generateUUID(),
		Name:    generateString(lowercaseChars, 10),
		StateID: generateUUID(),
	}
}

func GenerateCreateCityRequest() *payloads.CreateCityRequest {
	return &payloads.CreateCityRequest{
		Name: generateString(lowercaseChars, 10),
	}
}

// Helpers
const lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
const uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberChars = "0123456789"
const letterChars = lowercaseChars + uppercaseChars
const allChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

var urlProtocols = []string{"http", "https"}
var urlDomains = []string{"example.com", "testsite.com", "mywebsite.org", "randomsite.net"}
var emailDomains = []string{"gmail.com", "yahoo.com", "outlook.com", "example.com"}

func generateString(chars string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func generateBool() bool {
	return rand.Intn(2) == 1
}

func generateFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func generateInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func generateCurrentTime() time.Time {
	return time.Now().UTC()
}

func generateUUID() uuid.UUID {
	return uuid.New()
}
