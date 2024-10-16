package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestCountryService_GetCountries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockCountryRepository(ctrl)
	service := NewCountryService(nil, nil, repo)

	t.Run("success", func(t *testing.T) {
		countries := make([]*models.Country, 3)
		for i := 0; i < len(countries); i++ {
			countries[i] = utils.GenerateRandomCountry()
		}

		repo.EXPECT().GetCountries().Return(countries, nil).Times(1)

		result, err := service.GetCountries()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, countries, result)
	})

	t.Run("error getting countries", func(t *testing.T) {
		repo.EXPECT().GetCountries().Return(nil, errors.New("error getting countries")).Times(1)

		result, err := service.GetCountries()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting countries", err.Error())
	})
}

func TestCountryService_CreateCountry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockCountryRepository(ctrl)
	service := NewCountryService(nil, transaction, repo)

	country := utils.GenerateRandomCountry()

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetCountryByName(country.Name).Return(nil, nil).Times(1)
		repo.EXPECT().GetCountryByCode(country.Code).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateCountry(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateCountry(country.Name, country.Code)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, country.Name, result.Name)
		assert.Equal(t, country.Code, result.Code)
	})

	t.Run("duplicate name", func(t *testing.T) {
		repo.EXPECT().GetCountryByName(country.Name).Return(country, nil).Times(1)

		result, err := service.CreateCountry(country.Name, country.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate country name", err.Error())
	})

	t.Run("error getting by name", func(t *testing.T) {
		repo.EXPECT().GetCountryByName(country.Name).Return(nil, errors.New("error getting by name")).Times(1)

		result, err := service.CreateCountry(country.Name, country.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting by name", err.Error())
	})

	t.Run("duplicate code", func(t *testing.T) {
		repo.EXPECT().GetCountryByName(country.Name).Return(nil, nil).Times(1)
		repo.EXPECT().GetCountryByCode(country.Code).Return(country, nil).Times(1)

		result, err := service.CreateCountry(country.Name, country.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate country code", err.Error())
	})

	t.Run("error getting by code", func(t *testing.T) {
		repo.EXPECT().GetCountryByName(country.Name).Return(nil, nil).Times(1)
		repo.EXPECT().GetCountryByCode(country.Code).Return(nil, errors.New("error getting by code")).Times(1)

		result, err := service.CreateCountry(country.Name, country.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting by code", err.Error())
	})

	t.Run("error creating country", func(t *testing.T) {
		repo.EXPECT().GetCountryByName(country.Name).Return(nil, nil).Times(1)
		repo.EXPECT().GetCountryByCode(country.Code).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateCountry(gomock.Any(), gomock.Any()).Return(errors.New("error creating country")).Times(1)

		result, err := service.CreateCountry(country.Name, country.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating country", err.Error())
	})
}
