package services

import (
	"encoding/json"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserLocationService interface {
	GetCurrentUserLocation() (*models.UserLocation, *errors.ApiError)
}

func NewUserLocationService() UserLocationService {
	clientExpiresAfter, err := strconv.Atoi(os.Getenv("USER_LOCATION_API_TIMEOUT_SECONDS"))
	if err != nil {
		clientExpiresAfter = 10
	}
	httpClient := &http.Client{Timeout: time.Duration(clientExpiresAfter) * time.Second}

	return &userLocationService{
		url:        os.Getenv("USER_LOCATION_API_URL"),
		httpClient: httpClient,
	}
}

type userLocationService struct {
	url        string
	httpClient *http.Client
}

func (s *userLocationService) GetCurrentUserLocation() (*models.UserLocation, *errors.ApiError) {
	req, err := http.NewRequest(http.MethodGet, s.url, nil)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.InternalServerError("unexpected response from location API")
	}

	var location models.UserLocation
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &location, nil
}
