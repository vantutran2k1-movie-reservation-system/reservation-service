package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestCityService_CreateCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	countryRepo := mock_repositories.NewMockCountryRepository(ctrl)
	stateRepo := mock_repositories.NewMockStateRepository(ctrl)
	cityRepo := mock_repositories.NewMockCityRepository(ctrl)
	service := NewCityService(nil, transaction, countryRepo, stateRepo, cityRepo)

	country := utils.GenerateCountry()
	state := utils.GenerateState()
	city := utils.GenerateCity()

	t.Run("success", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(city.StateID).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCityByName(city.StateID, city.Name).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		cityRepo.EXPECT().CreateCity(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, city.Name)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, city.Name, result.Name)
	})

	t.Run("country not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(nil, nil).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, city.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "country does not exist", err.Error())
	})

	t.Run("error getting country", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(nil, errors.New("error getting country")).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, city.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting country", err.Error())
	})

	t.Run("state not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(city.StateID).Return(nil, nil).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, city.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "state does not exist", err.Error())
	})

	t.Run("error getting state", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(city.StateID).Return(nil, errors.New("error getting state")).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, city.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting state", err.Error())
	})

	t.Run("duplicate city name", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(city.StateID).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCityByName(city.StateID, city.Name).Return(city, nil).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, city.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate city name for this state", err.Error())
	})

	t.Run("error getting city", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(city.StateID).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCityByName(city.StateID, city.Name).Return(nil, errors.New("error getting city")).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, city.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting city", err.Error())
	})

	t.Run("error creating city", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetState(city.StateID).Return(state, nil).Times(1)
		cityRepo.EXPECT().GetCityByName(city.StateID, city.Name).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		cityRepo.EXPECT().CreateCity(gomock.Any(), gomock.Any()).Return(errors.New("error creating city")).Times(1)

		result, err := service.CreateCity(state.CountryID, city.StateID, city.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating city", err.Error())
	})
}
