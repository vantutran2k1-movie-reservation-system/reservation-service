package utils

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"golang.org/x/crypto/bcrypt"
)

func GetPointerOf[T any](value T) *T {
	return &value
}

// Models
func GenerateUser() *models.User {
	return &models.User{
		ID:           generateUUID(),
		Email:        generateEmail(),
		PasswordHash: generateHashedPassword(),
		IsActive:     generateBool(),
		IsVerified:   generateBool(),
		CreatedAt:    generateCurrentTime(),
		UpdatedAt:    generateCurrentTime(),
	}
}

func GenerateUserProfile() *models.UserProfile {
	return &models.UserProfile{
		ID:                generateUUID(),
		UserID:            generateUUID(),
		FirstName:         generateName(),
		LastName:          generateName(),
		PhoneNumber:       GetPointerOf(generatePhoneNumber()),
		DateOfBirth:       GetPointerOf(generateDate()),
		ProfilePictureUrl: GetPointerOf(generateURL()),
		Bio:               GetPointerOf(generateString(allChars, 50)),
		CreatedAt:         generateCurrentTime(),
		UpdatedAt:         generateCurrentTime(),
	}
}

func GenerateLoginToken() *models.LoginToken {
	return &models.LoginToken{
		ID:         generateUUID(),
		UserID:     generateUUID(),
		TokenValue: uuid.NewString(),
		CreatedAt:  generateCurrentTime(),
		ExpiresAt:  generateCurrentTime().Add(60 * time.Minute),
	}
}

func GenerateUserSession() *models.UserSession {
	return &models.UserSession{
		UserID: generateUUID(),
		Email:  generateEmail(),
	}
}

func GenerateMovie() *models.Movie {
	return &models.Movie{
		ID:              generateUUID(),
		Title:           generateString(letterChars, 10),
		Description:     GetPointerOf(generateString(letterChars, 100)),
		ReleaseDate:     generateDate(),
		DurationMinutes: generateInt(100, 200),
		Language:        GetPointerOf(generateString(letterChars, 10)),
		Rating:          GetPointerOf(generateFloat(0, 5)),
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

func GenerateUserRegistrationToken() *models.UserRegistrationToken {
	return &models.UserRegistrationToken{
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

func GenerateState() *models.State {
	return &models.State{
		ID:        generateUUID(),
		Name:      generateString(lowercaseChars, 10),
		Code:      GetPointerOf(generateString(uppercaseChars, 2)),
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
	return &models.TheaterLocation{
		ID:         generateUUID(),
		TheaterID:  GetPointerOf(generateUUID()),
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

func generateURL() string {
	protocol := urlProtocols[rand.Intn(len(urlProtocols))]
	subdomain := generateString(lowercaseChars, rand.Intn(5)+3)
	domain := urlDomains[rand.Intn(len(urlDomains))]
	path := generateString(lowercaseChars+"/", rand.Intn(10)+5)
	queryParams := generateString(lowercaseChars, rand.Intn(3)+3) + "=value"

	return fmt.Sprintf("%s://%s.%s/%s?%s", protocol, subdomain, domain, path, queryParams)
}
