package services

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUserLocationService_GetCurrentUserLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("USER_LOCATION_API_TIMEOUT_SECONDS", "10")
	defer os.Unsetenv("USER_LOCATION_API_TIMEOUT_SECONDS")

	os.Setenv("USER_LOCATION_API_URL", "https://example.com")
	defer os.Unsetenv("USER_LOCATION_API_URL")

	t.Run("success", func(t *testing.T) {
		mockResponse := `{"lat": 10.8231, "lon": 106.6297}`
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockResponse))
		}))
		defer mockServer.Close()

		os.Setenv("USER_LOCATION_API_URL", mockServer.URL)
		defer os.Unsetenv("USER_LOCATION_API_URL")

		service := NewUserLocationService()

		location, apiErr := service.GetCurrentUserLocation()

		assert.NotNil(t, location)
		assert.Nil(t, apiErr)
		assert.Equal(t, 10.8231, location.Latitude)
		assert.Equal(t, 106.6297, location.Longitude)
	})

	t.Run("API returns non-200 status code", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		os.Setenv("USER_LOCATION_API_URL", mockServer.URL)
		defer os.Unsetenv("USER_LOCATION_API_URL")

		service := NewUserLocationService()

		location, apiErr := service.GetCurrentUserLocation()

		assert.Nil(t, location)
		assert.NotNil(t, apiErr)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		assert.Equal(t, "unexpected response from location API", apiErr.Error())
	})

	t.Run("API response body is invalid JSON", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid JSON"))
		}))
		defer mockServer.Close()

		os.Setenv("USER_LOCATION_API_URL", mockServer.URL)
		defer os.Unsetenv("USER_LOCATION_API_URL")

		service := NewUserLocationService()

		location, apiErr := service.GetCurrentUserLocation()

		assert.Nil(t, location)
		assert.NotNil(t, apiErr)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		assert.Contains(t, apiErr.Error(), "invalid character")
	})
}
