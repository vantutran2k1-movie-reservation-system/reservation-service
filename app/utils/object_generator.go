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

func GenerateURL() string {
	protocol := urlProtocols[rand.Intn(len(urlProtocols))]
	subdomain := generateString(lowercaseChars, rand.Intn(5)+3)
	domain := urlDomains[rand.Intn(len(urlDomains))]
	path := generateString(lowercaseChars+"/", rand.Intn(10)+5)
	queryParams := generateString(lowercaseChars, rand.Intn(3)+3) + "=value"

	return fmt.Sprintf("%s://%s.%s/%s?%s", protocol, subdomain, domain, path, queryParams)
}

// Models and payloads
//
// User
func GenerateUser() *models.User {
	return &models.User{
		ID:           generateUUID(),
		Email:        generateEmail(),
		PasswordHash: generateHashedPassword(),
		CreatedAt:    generateCurrentTime(),
		UpdatedAt:    generateCurrentTime(),
	}
}

func GenerateCreateUserRequest() payloads.CreateUserRequest {
	return payloads.CreateUserRequest{
		Email:    generateEmail(),
		Password: generatePassword(),
	}
}

func GenerateUpdatePasswordRequest() payloads.UpdatePasswordRequest {
	return payloads.UpdatePasswordRequest{
		Password: generatePassword(),
	}
}

func GenerateLoginUserRequest() payloads.LoginUserRequest {
	return payloads.LoginUserRequest{
		Email:    generateEmail(),
		Password: generatePassword(),
	}
}

func GenerateResetUserPasswordRequest() payloads.ResetPasswordRequest {
	return payloads.ResetPasswordRequest{
		Password: generatePassword(),
	}
}

// User profile
func GenerateUserProfile() *models.UserProfile {
	phoneNumber := generatePhoneNumber()
	dateOfBirth := generateDate()
	profilePictureUrl := GenerateURL()
	bio := generateString(allChars, 50)

	return &models.UserProfile{
		ID:                generateUUID(),
		UserID:            generateUUID(),
		FirstName:         generateName(),
		LastName:          generateName(),
		PhoneNumber:       &phoneNumber,
		DateOfBirth:       &dateOfBirth,
		ProfilePictureUrl: &profilePictureUrl,
		Bio:               &bio,
		CreatedAt:         generateCurrentTime(),
		UpdatedAt:         generateCurrentTime(),
	}
}

func GenerateCreateUserProfileRequest() payloads.CreateUserProfileRequest {
	phoneNumber := generatePhoneNumber()
	dateOfBirth := generateDate()

	return payloads.CreateUserProfileRequest{
		FirstName:   generateName(),
		LastName:    generateName(),
		PhoneNumber: &phoneNumber,
		DateOfBirth: &dateOfBirth,
	}
}

func GenerateUpdateUserProfileRequest() payloads.UpdateUserProfileRequest {
	phoneNumber := generatePhoneNumber()
	dateOfBirth := generateDate()

	return payloads.UpdateUserProfileRequest{
		FirstName:   generateName(),
		LastName:    generateName(),
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
		Email:  generateEmail(),
	}
}

func GenerateSessionID() string {
	return uuid.NewString()
}

// File
func GenerateFileHeader() *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: generateName(),
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
		ReleaseDate:     generateDate(),
		DurationMinutes: generateInt(100, 200),
		Language:        &language,
		Rating:          &rating,
		IsActive:        generateBool(),
		CreatedAt:       generateCurrentTime(),
		UpdatedAt:       generateCurrentTime(),
		IsDeleted:       generateBool(),
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

func GenerateCreateMovieRequest() payloads.CreateMovieRequest {
	description := generateString(allChars, 100)
	language := generateString(lowercaseChars, 10)
	rating := generateFloat(0, 5)
	isActive := generateBool()

	return payloads.CreateMovieRequest{
		Title:           generateString(allChars, 10),
		Description:     &description,
		ReleaseDate:     generateDate(),
		DurationMinutes: generateInt(100, 200),
		Language:        &language,
		Rating:          &rating,
		IsActive:        &isActive,
	}
}

func GenerateUpdateMovieRequest() payloads.UpdateMovieRequest {
	description := generateString(allChars, 100)
	language := generateString(lowercaseChars, 10)
	rating := generateFloat(0, 5)
	isActive := generateBool()

	return payloads.UpdateMovieRequest{
		Title:           generateString(allChars, 10),
		Description:     &description,
		ReleaseDate:     generateDate(),
		DurationMinutes: generateInt(100, 200),
		Language:        &language,
		Rating:          &rating,
		IsActive:        &isActive,
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

func GenerateCreateGenreRequest() payloads.CreateGenreRequest {
	return payloads.CreateGenreRequest{
		Name: generateString(letterChars, 10),
	}
}

func GenerateUpdateGenreRequest() payloads.UpdateGenreRequest {
	return payloads.UpdateGenreRequest{
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

func GenerateCreatePasswordResetTokenRequest() payloads.CreatePasswordResetTokenRequest {
	return payloads.CreatePasswordResetTokenRequest{
		Email: generateEmail(),
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

func GenerateCreateCountryRequest() payloads.CreateCountryRequest {
	return payloads.CreateCountryRequest{
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

func GenerateCreateStateRequest() payloads.CreateStateRequest {
	code := generateString(uppercaseChars, 2)

	return payloads.CreateStateRequest{
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

func GenerateCities(count int) []*models.City {
	cities := make([]*models.City, count)
	for i := 0; i < count; i++ {
		cities[i] = GenerateCity()
	}

	return cities
}

func GenerateCreateCityRequest() payloads.CreateCityRequest {
	return payloads.CreateCityRequest{
		Name: generateString(lowercaseChars, 10),
	}
}

// Theater
func GenerateTheater() *models.Theater {
	return &models.Theater{
		ID:   generateUUID(),
		Name: generateString(letterChars, 10),
	}
}

func GenerateTheaters(count int) []*models.Theater {
	theaters := make([]*models.Theater, count)
	for i := 0; i < count; i++ {
		theaters[i] = GenerateTheater()
	}

	return theaters
}

func GenerateTheaterLocation() *models.TheaterLocation {
	theaterID := generateUUID()
	return &models.TheaterLocation{
		ID:         generateUUID(),
		TheaterID:  &theaterID,
		CityID:     generateUUID(),
		Address:    generateString(lowercaseChars, 50),
		PostalCode: generateString(numberChars, 6),
		Latitude:   generateFloat(1, 100),
		Longitude:  generateFloat(1, 100),
	}
}

func GenerateTheaterLocations(count int) []*models.TheaterLocation {
	locations := make([]*models.TheaterLocation, count)
	for i := 0; i < count; i++ {
		locations[i] = GenerateTheaterLocation()
	}

	return locations
}

func GenerateUserLocation() *models.UserLocation {
	return &models.UserLocation{
		Latitude:  generateFloat(1, 100),
		Longitude: generateFloat(1, 100),
	}
}

func GenerateCreateTheaterRequest() payloads.CreateTheaterRequest {
	return payloads.CreateTheaterRequest{
		Name: generateString(letterChars, 10),
	}
}

func GenerateCreateTheaterLocationRequest() payloads.CreateTheaterLocationRequest {
	return payloads.CreateTheaterLocationRequest{
		CityID:     generateUUID(),
		Address:    generateString(lowercaseChars, 50),
		PostalCode: generateString(numberChars, 6),
		Latitude:   generateFloat(1.0, 100.0),
		Longitude:  generateFloat(1.0, 100.0),
	}
}

func GenerateUpdateTheaterLocationRequest() payloads.UpdateTheaterLocationRequest {
	return payloads.UpdateTheaterLocationRequest{
		CityID:     generateUUID(),
		Address:    generateString(lowercaseChars, 50),
		PostalCode: generateString(numberChars, 6),
		Latitude:   generateFloat(1.0, 100.0),
		Longitude:  generateFloat(1.0, 100.0),
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

func generateEmail() string {
	nameLength := rand.Intn(6) + 5
	domain := emailDomains[rand.Intn(len(emailDomains))]

	return fmt.Sprintf("%s@%s", generateString(lowercaseChars, nameLength), domain)
}

func generatePassword() string {
	return generateString(letterChars, 12)
}

func generateHashedPassword() string {
	password := generatePassword()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword)
}

func generatePhoneNumber() string {
	return generateString(numberChars, 10)
}

func generateName() string {
	return generateString(lowercaseChars, 5)
}

func generateDate() string {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := generateCurrentTime()

	days := rand.Int63n(int64(end.Sub(start).Hours() / 24))

	date := start.AddDate(0, 0, int(days)-1)

	return date.Format("2006-01-02")
}
