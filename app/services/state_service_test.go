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

func TestStateService_CreateState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	countryRepo := mock_repositories.NewMockCountryRepository(ctrl)
	stateRepo := mock_repositories.NewMockStateRepository(ctrl)
	service := NewStateService(nil, transaction, countryRepo, stateRepo)

	country := utils.GenerateRandomCountry()
	state := utils.GenerateRandomState()

	t.Run("success", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetStateByName(state.CountryID, state.Name).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		stateRepo.EXPECT().CreateState(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateState(state.CountryID, state.Name, state.Code)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, state.Name, result.Name)
		assert.Equal(t, state.Code, result.Code)
		assert.Equal(t, state.CountryID, result.CountryID)
	})

	t.Run("country not found", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(nil, nil).Times(1)

		result, err := service.CreateState(state.CountryID, state.Name, state.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "country does not exist", err.Error())
	})

	t.Run("error getting country", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(nil, errors.New("error getting country")).Times(1)

		result, err := service.CreateState(state.CountryID, state.Name, state.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting country", err.Error())
	})

	t.Run("duplicate state name", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetStateByName(state.CountryID, state.Name).Return(state, nil).Times(1)

		result, err := service.CreateState(state.CountryID, state.Name, state.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate state name for this country", err.Error())
	})

	t.Run("error getting state", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetStateByName(state.CountryID, state.Name).Return(nil, errors.New("error getting state")).Times(1)

		result, err := service.CreateState(state.CountryID, state.Name, state.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting state", err.Error())
	})

	t.Run("error creating state", func(t *testing.T) {
		countryRepo.EXPECT().GetCountry(state.CountryID).Return(country, nil).Times(1)
		stateRepo.EXPECT().GetStateByName(state.CountryID, state.Name).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		stateRepo.EXPECT().CreateState(gomock.Any(), gomock.Any()).Return(errors.New("error creating state")).Times(1)

		result, err := service.CreateState(state.CountryID, state.Name, state.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating state", err.Error())
	})
}
