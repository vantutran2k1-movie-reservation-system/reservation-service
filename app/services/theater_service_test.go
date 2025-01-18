package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	apiError "github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_transaction"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestTheaterService_GetTheater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockTheaterRepository(ctrl)
	service := NewTheaterService(nil, nil, repo, nil, nil, nil, nil)

	theater := utils.GenerateTheater()
	filter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: theater.ID},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(theater, nil).Times(1)

		result, err := service.GetTheater(theater.ID, false)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, theater, result)
	})

	t.Run("theater not found", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, nil).Times(1)

		result, err := service.GetTheater(theater.ID, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "theater not found", err.Error())
	})

	t.Run("error getting theater", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.GetTheater(theater.ID, false)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theater", err.Error())
	})
}

func TestTheaterService_GetTheaters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockTheaterRepository(ctrl)
	service := NewTheaterService(nil, nil, repo, nil, nil, nil, nil)

	theaters := utils.GenerateTheaters(3)

	limit := 2
	offset := 0
	includeLocation := true
	getFilter := filters.TheaterFilter{
		Filter: &filters.MultiFilter{
			Limit:  &limit,
			Offset: &offset,
		},
	}
	countFilter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetTheaters(getFilter, includeLocation).Return(theaters, nil).Times(1)
		repo.EXPECT().GetNumbersOfTheater(countFilter).Return(len(theaters), nil).Times(1)

		result, meta, err := service.GetTheaters(limit, offset, includeLocation)

		assert.NotNil(t, result)
		assert.NotNil(t, meta)
		assert.Nil(t, err)
		assert.Equal(t, theaters, result)

		nextUrl := fmt.Sprintf("/theaters?%s=%d&%s=%d&%s=%v", constants.Limit, limit, constants.Offset, offset+limit, constants.IncludeTheaterLocation, includeLocation)
		expectedMeta := models.ResponseMeta{
			Limit:   limit,
			Offset:  offset,
			Total:   len(theaters),
			NextUrl: &nextUrl,
		}
		assert.Equal(t, &expectedMeta, meta)
	})

	t.Run("error getting theaters", func(t *testing.T) {
		repo.EXPECT().GetTheaters(getFilter, includeLocation).Return(nil, errors.New("error getting theaters")).Times(1)

		result, meta, err := service.GetTheaters(limit, offset, includeLocation)

		assert.Nil(t, result)
		assert.Nil(t, meta)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theaters", err.Error())
	})

	t.Run("error counting theaters", func(t *testing.T) {
		repo.EXPECT().GetTheaters(getFilter, includeLocation).Return(theaters, nil).Times(1)
		repo.EXPECT().GetNumbersOfTheater(countFilter).Return(0, errors.New("error counting theaters")).Times(1)

		result, meta, err := service.GetTheaters(limit, offset, includeLocation)

		assert.Nil(t, result)
		assert.Nil(t, meta)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error counting theaters", err.Error())
	})
}

func TestTheaterService_GetNearbyTheaters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repositories.NewMockTheaterRepository(ctrl)
	userLocService := mock_services.NewMockUserLocationService(ctrl)
	service := NewTheaterService(nil, nil, repo, nil, nil, nil, userLocService)

	userLoc := &models.UserLocation{
		Latitude:  20.0,
		Longitude: 30.0,
	}
	distance := 10.0

	t.Run("success", func(t *testing.T) {
		getTheaterResults := make([]*payloads.GetTheaterWithLocationResult, 3)
		for i := 0; i < len(getTheaterResults); i++ {
			theater := utils.GenerateTheater()
			location := utils.GenerateTheaterLocation()
			getTheaterResults[i] = &payloads.GetTheaterWithLocationResult{
				Id:         theater.ID,
				Name:       theater.Name,
				LocationId: location.ID,
				CityId:     location.CityID,
				Address:    location.Address,
				PostalCode: location.PostalCode,
				Latitude:   location.Latitude,
				Longitude:  location.Longitude,
			}
		}

		userLocService.EXPECT().GetCurrentUserLocation().Return(userLoc, nil).Times(1)
		repo.EXPECT().GetNearbyTheatersWithLocations(userLoc.Latitude, userLoc.Longitude, distance).Return(getTheaterResults, nil).Times(1)

		theaters, err := service.GetNearbyTheaters(distance)

		assert.NotNil(t, theaters)
		assert.Nil(t, err)
		for i, theater := range theaters {
			assert.Equal(t, getTheaterResults[i].Id, theater.ID)
			assert.Equal(t, getTheaterResults[i].LocationId, theater.Location.ID)
		}
	})

	t.Run("error getting current user location", func(t *testing.T) {
		userLocService.EXPECT().GetCurrentUserLocation().Return(nil, apiError.InternalServerError("error getting current user location")).Times(1)

		theaters, err := service.GetNearbyTheaters(distance)

		assert.Nil(t, theaters)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting current user location", err.Error())
	})

	t.Run("error getting theaters", func(t *testing.T) {
		userLocService.EXPECT().GetCurrentUserLocation().Return(userLoc, nil).Times(1)
		repo.EXPECT().GetNearbyTheatersWithLocations(userLoc.Latitude, userLoc.Longitude, distance).Return(nil, errors.New("error getting theaters")).Times(1)

		theaters, err := service.GetNearbyTheaters(distance)

		assert.Nil(t, theaters)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theaters", err.Error())
	})
}

func TestTheaterService_CreateTheater(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	repo := mock_repositories.NewMockTheaterRepository(ctrl)
	service := NewTheaterService(nil, transaction, repo, nil, nil, nil, nil)

	theater := utils.GenerateTheater()
	req := payloads.CreateTheaterRequest{
		Name: theater.Name,
	}
	filter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		Name:   &filters.Condition{Operator: filters.OpEqual, Value: &req.Name},
	}

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateTheater(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateTheater(req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.Name, result.Name)
	})

	t.Run("duplicate theater name", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(theater, nil).Times(1)

		result, err := service.CreateTheater(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate theater name", err.Error())
	})

	t.Run("error getting theater", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.CreateTheater(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theater", err.Error())
	})

	t.Run("error creating theater", func(t *testing.T) {
		repo.EXPECT().GetTheater(filter, false).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		repo.EXPECT().CreateTheater(gomock.Any(), gomock.Any()).Return(errors.New("error creating theater")).Times(1)

		result, err := service.CreateTheater(req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating theater", err.Error())
	})
}

func TestTheaterService_CreateTheaterLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	theaterRepo := mock_repositories.NewMockTheaterRepository(ctrl)
	theaterLocationRepo := mock_repositories.NewMockTheaterLocationRepository(ctrl)
	cityRepo := mock_repositories.NewMockCityRepository(ctrl)
	service := NewTheaterService(nil, transaction, theaterRepo, theaterLocationRepo, nil, cityRepo, nil)

	theater := utils.GenerateTheater()
	city := utils.GenerateCity()
	req := payloads.CreateTheaterLocationRequest{
		CityID:     uuid.New(),
		Address:    "address",
		PostalCode: "700000",
		Latitude:   50.0,
		Longitude:  100.0,
	}
	cityFilter := filters.CityFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: req.CityID},
	}
	theaterFilter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: theater.ID},
	}

	t.Run("success", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(city, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		theaterLocationRepo.EXPECT().CreateTheaterLocation(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, req.CityID, result.CityID)
		assert.Equal(t, req.Address, result.Address)
		assert.Equal(t, req.PostalCode, result.PostalCode)
		assert.Equal(t, req.Latitude, result.Latitude)
		assert.Equal(t, req.Longitude, result.Longitude)
	})

	t.Run("theater not found", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(nil, nil).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "theater not found", err.Error())
	})

	t.Run("error getting theater", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theater", err.Error())
	})

	t.Run("duplicate theater location", func(t *testing.T) {
		th := utils.GenerateTheater()
		th.Location = utils.GenerateTheaterLocation()
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(th, nil).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "duplicate location for this theater", err.Error())
	})

	t.Run("city not found", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(nil, nil).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "invalid city id", err.Error())
	})

	t.Run("error getting city", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(nil, errors.New("error getting city")).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting city", err.Error())
	})

	t.Run("error creating location", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(city, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		theaterLocationRepo.EXPECT().CreateTheaterLocation(gomock.Any(), gomock.Any()).Return(errors.New("error creating location")).Times(1)

		result, err := service.CreateTheaterLocation(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error creating location", err.Error())
	})
}

func TestTheaterService_CreateSeat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	theaterRepo := mock_repositories.NewMockTheaterRepository(ctrl)
	seatRepo := mock_repositories.NewMockSeatRepository(ctrl)

	service := NewTheaterService(nil, transaction, theaterRepo, nil, seatRepo, nil, nil)

	theater := utils.GenerateTheater()
	seat := utils.GenerateSeat()
	seat.TheaterId = &theater.ID
	req := payloads.CreateSeatPayload{
		Row:    seat.Row,
		Number: seat.Number,
		Type:   seat.Type,
	}
	theaterFilter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: theater.ID},
	}
	seatFilter := filters.SeatFilter{
		Filter:    &filters.SingleFilter{},
		TheaterId: &filters.Condition{Operator: filters.OpEqual, Value: theater.ID},
		Row:       &filters.Condition{Operator: filters.OpEqual, Value: req.Row},
		Number:    &filters.Condition{Operator: filters.OpEqual, Value: req.Number},
	}

	t.Run("success", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(theater, nil).Times(1)
		seatRepo.EXPECT().GetSeat(seatFilter).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		seatRepo.EXPECT().CreateSeat(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := service.CreateSeat(theater.ID, req)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, seat.TheaterId, result.TheaterId)
		assert.Equal(t, seat.Row, result.Row)
		assert.Equal(t, seat.Number, result.Number)
		assert.Equal(t, seat.Type, result.Type)
	})

	t.Run("theater not found", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(nil, nil).Times(1)

		result, err := service.CreateSeat(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.EqualError(t, err, "theater not found")
	})

	t.Run("error getting theater", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(nil, errors.New("error getting theater")).Times(1)

		result, err := service.CreateSeat(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.EqualError(t, err, "error getting theater")
	})

	t.Run("duplicate seat", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(theater, nil).Times(1)
		seatRepo.EXPECT().GetSeat(seatFilter).Return(seat, nil).Times(1)

		result, err := service.CreateSeat(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.EqualError(t, err, "duplicate seat for this theater")
	})

	t.Run("error getting seat", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(theater, nil).Times(1)
		seatRepo.EXPECT().GetSeat(seatFilter).Return(nil, errors.New("error getting seat")).Times(1)

		result, err := service.CreateSeat(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.EqualError(t, err, "error getting seat")
	})

	t.Run("error creating seat", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, false).Return(theater, nil).Times(1)
		seatRepo.EXPECT().GetSeat(seatFilter).Return(nil, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		seatRepo.EXPECT().CreateSeat(gomock.Any(), gomock.Any()).Return(errors.New("error creating seat")).Times(1)

		result, err := service.CreateSeat(theater.ID, req)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.EqualError(t, err, "error creating seat")
	})
}

func TestTheaterService_UpdateTheaterLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transaction := mock_transaction.NewMockTransactionManager(ctrl)
	theaterRepo := mock_repositories.NewMockTheaterRepository(ctrl)
	theaterLocationRepo := mock_repositories.NewMockTheaterLocationRepository(ctrl)
	cityRepo := mock_repositories.NewMockCityRepository(ctrl)
	service := NewTheaterService(nil, transaction, theaterRepo, theaterLocationRepo, nil, cityRepo, nil)

	theater := utils.GenerateTheater()
	location := utils.GenerateTheaterLocation()
	location.TheaterID = &theater.ID
	theater.Location = location
	city := utils.GenerateCity()
	req := payloads.UpdateTheaterLocationRequest{
		CityID:     uuid.New(),
		Address:    "address",
		PostalCode: "700000",
		Latitude:   50.0,
		Longitude:  100.0,
	}
	cityFilter := filters.CityFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: req.CityID},
	}
	theaterFilter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: theater.ID},
	}

	t.Run("success", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(city, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		theaterLocationRepo.EXPECT().UpdateTheaterLocation(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		l, err := service.UpdateTheaterLocation(theater.ID, req)

		assert.NotNil(t, l)
		assert.Nil(t, err)
		assert.Equal(t, req.CityID, l.CityID)
		assert.Equal(t, req.Address, l.Address)
		assert.Equal(t, req.PostalCode, l.PostalCode)
		assert.Equal(t, req.Latitude, l.Latitude)
		assert.Equal(t, req.Longitude, l.Longitude)
	})

	t.Run("theater not found", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(nil, nil).Times(1)

		l, err := service.UpdateTheaterLocation(theater.ID, req)

		assert.Nil(t, l)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
		assert.Equal(t, "theater not found", err.Error())
	})

	t.Run("error getting theater", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(nil, errors.New("error getting theater")).Times(1)

		l, err := service.UpdateTheaterLocation(theater.ID, req)

		assert.Nil(t, l)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error getting theater", err.Error())
	})

	t.Run("theater location not found", func(t *testing.T) {
		th := utils.GenerateTheater()
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(th, nil).Times(1)

		l, err := service.UpdateTheaterLocation(theater.ID, req)

		assert.Nil(t, l)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, "location not found", err.Error())
	})

	t.Run("error updating location", func(t *testing.T) {
		theaterRepo.EXPECT().GetTheater(theaterFilter, true).Return(theater, nil).Times(1)
		cityRepo.EXPECT().GetCity(cityFilter).Return(city, nil).Times(1)
		transaction.EXPECT().ExecuteInTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(db *gorm.DB, fn func(tx *gorm.DB) error) error {
				return fn(db)
			},
		).Times(1)
		theaterLocationRepo.EXPECT().UpdateTheaterLocation(gomock.Any(), gomock.Any()).Return(errors.New("error updating location")).Times(1)

		l, err := service.UpdateTheaterLocation(theater.ID, req)

		assert.Nil(t, l)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
		assert.Equal(t, "error updating location", err.Error())
	})
}
