package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestLocationService_GetCountries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockCountryRepository(ctrl)
	service := NewLocationService(nil, nil, repo, nil, nil)

	t.Run("success", func(t *testing.T) {
		countries := utils.GenerateCountries(3)

		repo.EXPECT().GetCountries(gomock.Any()).Return(countries, nil).Times(1)

		result, err := service.GetCountries()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, countries, result)
	})

	t.Run("error getting countries", func(t *testing.T) {
		repo.EXPECT().GetCountries(gomock.Any()).Return(nil, errors.New("error getting countries")).Times(1)

		result, err := service.GetCountries()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting countries", err.Error())
	})
}

func TestLocationService_CreateCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockCountryRepository(ctrl)
	service := NewLocationService(nil, transaction, repo, nil, nil)

	country := utils.GenerateCountry()
	req := payloads.CreateCountryRequest{
		Name: country.Name,
		Code: country.Code,
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetCountry(gomock.Any()).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateCountry(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateCountry(req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.Code, result.Code)
	})

	t.Run("duplicate name or code", func(t *testing.T) {
		repo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)

		result, err := service.CreateCountry(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate country name or code", err.Error())
	})

	t.Run("error getting country", func(t *testing.T) {
		repo.EXPECT().GetCountry(gomock.Any()).Return(nil, errors.New("error getting country")).Times(1)

		result, err := service.CreateCountry(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting country", err.Error())
	})

	t.Run("error creating country", func(t *testing.T) {
		repo.EXPECT().GetCountry(gomock.Any()).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateCountry(gomock.Any(), gomock.Any()).Return(errors.New("error creating country")).Times(1)

		result, err := service.CreateCountry(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating country", err.Error())
	})
}

func TestLocationService_GetStatesByCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	countryRepo := mock_repositories.NewMockCountryRepository(ctrl)
	stateRepo := mock_repositories.NewMockStateRepository(ctrl)
	service := NewLocationService(nil, nil, countryRepo, stateRepo, nil)

	country := utils.GenerateCountry()

	t.Run("success", func(t *testing.T) {
		states := utils.GenerateStates(3)

		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetStates(gomock.Any()).Return(states, nil).Times(1)

		result, err := service.GetStatesByCountry(country.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, states, result)
	})

	t.Run("country not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(nil, nil).Times(1)

		result, err := service.GetStatesByCountry(country.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "country does not exist", err.Error())
	})

	t.Run("error getting country", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(nil, errors.New("error getting country")).Times(1)

		result, err := service.GetStatesByCountry(country.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting country", err.Error())
	})

	t.Run("error getting states", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetStates(gomock.Any()).Return(nil, errors.New("error getting states")).Times(1)

		result, err := service.GetStatesByCountry(country.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting states", err.Error())
	})
}

func TestLocationService_CreateState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	countryRepo := mock_repositories.NewMockCountryRepository(ctrl)
	stateRepo := mock_repositories.NewMockStateRepository(ctrl)
	service := NewLocationService(nil, transaction, countryRepo, stateRepo, nil)

	country := utils.GenerateCountry()
	state := utils.GenerateState()
	req := payloads.CreateStateRequest{
		Name: state.Name,
		Code: state.Code,
	}

	t.Run("success", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		stateRepo.EXPECT().CreateState(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateState(state.CountryID, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.Code, result.Code)
		assert.Equal(t, state.CountryID, result.CountryID)
	})

	t.Run("country not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(nil, nil).Times(1)

		result, err := service.CreateState(state.CountryID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "country does not exist", err.Error())
	})

	t.Run("error getting country", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(nil, errors.New("error getting country")).Times(1)

		result, err := service.CreateState(state.CountryID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting country", err.Error())
	})

	t.Run("duplicate state name", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(state, nil).Times(1)

		result, err := service.CreateState(state.CountryID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate state name for this country", err.Error())
	})

	t.Run("error getting state", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(nil, errors.New("error getting state")).Times(1)

		result, err := service.CreateState(state.CountryID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting state", err.Error())
	})

	t.Run("error creating state", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		stateRepo.EXPECT().CreateState(gomock.Any(), gomock.Any()).Return(errors.New("error creating state")).Times(1)

		result, err := service.CreateState(state.CountryID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating state", err.Error())
	})
}

func TestLocationService_GetCitiesByState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	countryRepo := mock_repositories.NewMockCountryRepository(ctrl)
	stateRepo := mock_repositories.NewMockStateRepository(ctrl)
	cityRepo := mock_repositories.NewMockCityRepository(ctrl)
	service := NewLocationService(nil, nil, countryRepo, stateRepo, cityRepo)

	country := utils.GenerateCountry()
	state := utils.GenerateState()
	cities := utils.GenerateCities(3)

	t.Run("success", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCities(gomock.Any()).Return(cities, nil).Times(1)

		result, err := service.GetCitiesByState(country.ID, state.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, cities, result)
	})

	t.Run("country not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(nil, nil).Times(1)

		result, err := service.GetCitiesByState(country.ID, state.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "country does not exist", err.Error())
	})

	t.Run("error getting country", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(nil, errors.New("error getting country")).Times(1)

		result, err := service.GetCitiesByState(country.ID, state.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting country", err.Error())
	})

	t.Run("state not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(nil, nil).Times(1)

		result, err := service.GetCitiesByState(state.CountryID, state.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "state does not exist", err.Error())
	})

	t.Run("error getting state", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(nil, errors.New("error getting state")).Times(1)

		result, err := service.GetCitiesByState(country.ID, state.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting state", err.Error())
	})

	t.Run("error getting cities", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCities(gomock.Any()).Return(nil, errors.New("error getting cities")).Times(1)

		result, err := service.GetCitiesByState(country.ID, state.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting cities", err.Error())
	})
}

func TestLocationService_CreateCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	countryRepo := mock_repositories.NewMockCountryRepository(ctrl)
	stateRepo := mock_repositories.NewMockStateRepository(ctrl)
	cityRepo := mock_repositories.NewMockCityRepository(ctrl)
	service := NewLocationService(nil, transaction, countryRepo, stateRepo, cityRepo)

	country := utils.GenerateCountry()
	state := utils.GenerateState()
	city := utils.GenerateCity()
	req := payloads.CreateCityRequest{
		Name: city.Name,
	}

	t.Run("success", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCity(gomock.Any()).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		cityRepo.EXPECT().CreateCity(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.Name, result.Name)
	})

	t.Run("country not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(nil, nil).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "country does not exist", err.Error())
	})

	t.Run("error getting country", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(nil, errors.New("error getting country")).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting country", err.Error())
	})

	t.Run("state not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(nil, nil).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "state does not exist", err.Error())
	})

	t.Run("error getting state", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(nil, errors.New("error getting state")).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting state", err.Error())
	})

	t.Run("duplicate city name", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCity(gomock.Any()).Return(city, nil).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate city name for this state", err.Error())
	})

	t.Run("error getting city", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCity(gomock.Any()).Return(nil, errors.New("error getting city")).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting city", err.Error())
	})

	t.Run("error creating city", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(gomock.Any()).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(gomock.Any()).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCity(gomock.Any()).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		cityRepo.EXPECT().CreateCity(gomock.Any(), gomock.Any()).Return(errors.New("error creating city")).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating city", err.Error())
	})
}
