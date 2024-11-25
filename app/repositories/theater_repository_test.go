package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
)

func TestTheaterRepository_GetTheater(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterRepository(db)

	theater := utils.GenerateTheater()
	location := utils.GenerateTheaterLocation()
	location.TheaterID = &theater.ID
	filter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: theater.ID},
		Name:   &filters.Condition{Operator: filters.OpEqual, Value: theater.Name},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" WHERE id = $1 AND name = $2 ORDER BY "theaters"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Name.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(theater))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theater_locations" WHERE "theater_locations"."theater_id" = $1`)).
			WithArgs(theater.ID).
			WillReturnRows(utils.GenerateSqlMockRow(location))

		result, err := repo.GetTheater(filter, true)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, location, result.Location)
	})

	t.Run("theater not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" WHERE id = $1 AND name = $2 ORDER BY "theaters"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Name.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetTheater(filter, true)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting theater", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" WHERE id = $1 AND name = $2 ORDER BY "theaters"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Name.Value, 1).
			WillReturnError(errors.New("error getting theater"))

		result, err := repo.GetTheater(filter, true)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting theater", err.Error())
	})
}

func TestTheaterRepository_GetTheaters(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterRepository(db)

	theaters := utils.GenerateTheaters(3)
	locations := utils.GenerateTheaterLocations(3)
	for i := range theaters {
		theaters[i].Location = locations[i]
		locations[i].TheaterID = &theaters[i].ID
	}
	limit := 3
	offset := 1
	filter := filters.TheaterFilter{
		Filter: &filters.MultiFilter{Limit: &limit, Offset: &offset, Sort: []filters.SortOption{{Field: "name", Direction: filters.Desc}}},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" ORDER BY name DESC LIMIT $1 OFFSET $2`)).
			WithArgs(limit, offset).
			WillReturnRows(utils.GenerateSqlMockRows(theaters))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theater_locations" WHERE "theater_locations"."theater_id" IN ($1,$2,$3)`)).
			WithArgs(theaters[0].ID, theaters[1].ID, theaters[2].ID).
			WillReturnRows(utils.GenerateSqlMockRows(locations))

		result, err := repo.GetTheaters(filter, true)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, theaters, result)
	})

	t.Run("error getting theaters", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" ORDER BY name DESC LIMIT $1 OFFSET $2`)).
			WithArgs(limit, offset).
			WillReturnError(errors.New("error getting theaters"))

		result, err := repo.GetTheaters(filter, true)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting theaters", err.Error())
	})

	t.Run("error getting locations", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" ORDER BY name DESC LIMIT $1 OFFSET $2`)).
			WithArgs(limit, offset).
			WillReturnRows(utils.GenerateSqlMockRows(theaters))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theater_locations" WHERE "theater_locations"."theater_id" IN ($1,$2,$3)`)).
			WithArgs(theaters[0].ID, theaters[1].ID, theaters[2].ID).
			WillReturnError(errors.New("error getting locations"))

		result, err := repo.GetTheaters(filter, true)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting locations", err.Error())
	})
}

func TestTheaterRepository_GetNearbyTheatersWithLocations(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterRepository(db)

	lat := 100.0
	lon := 101.0
	distance := 10.0
	expectedQuery := regexp.QuoteMeta(`
		WITH data AS (
			SELECT
				t.id, t.name,
				tl.id AS location_id, tl.city_id, tl.address, tl.postal_code, tl.latitude, tl.longitude,
				(6371 * acos(cos(radians($1)) * cos(radians(tl.latitude)) * cos(radians(tl.longitude) - radians($2)) + sin(radians($3)) * sin(radians(tl.latitude)))) AS distance
			FROM theaters t
			JOIN theater_locations tl ON t.id = tl.theater_id
		)
		SELECT *
		FROM data
		WHERE distance < $4
		ORDER BY distance
	`)

	t.Run("success", func(t *testing.T) {
		expectedResult := make([]*payloads.GetTheaterWithLocationResult, 3)
		for i := 0; i < len(expectedResult); i++ {
			theater := utils.GenerateTheater()
			location := utils.GenerateTheaterLocation()
			expectedResult[i] = &payloads.GetTheaterWithLocationResult{
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

		mock.ExpectQuery(expectedQuery).
			WithArgs(lat, lon, lat, distance).
			WillReturnRows(utils.GenerateSqlMockRows(expectedResult))

		result, err := repo.GetNearbyTheatersWithLocations(lat, lon, distance)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("error getting theaters", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).
			WithArgs(lat, lon, lat, distance).
			WillReturnError(errors.New("error getting theaters"))

		result, err := repo.GetNearbyTheatersWithLocations(lat, lon, distance)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting theaters", err.Error())
	})
}

func TestTheaterRepository_GetNumbersOfTheater(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterRepository(db)

	filter := filters.TheaterFilter{
		Filter: &filters.SingleFilter{},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "theaters"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		result, err := repo.GetNumbersOfTheater(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, 3, result)
	})

	t.Run("error counting theaters", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "theaters"`)).
			WillReturnError(errors.New("error counting theaters"))

		result, err := repo.GetNumbersOfTheater(filter)

		assert.NotNil(t, err)
		assert.Equal(t, 0, result)
		assert.Equal(t, "error counting theaters", err.Error())
	})
}

func TestTheaterRepository_CreateTheater(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterRepository(db)

	theater := utils.GenerateTheater()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "theaters" ("id","name") VALUES ($1,$2)`)).
			WithArgs(theater.ID, theater.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateTheater(tx, theater)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error creating theater", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "theaters" ("id","name") VALUES ($1,$2)`)).
			WithArgs(theater.ID, theater.Name).
			WillReturnError(errors.New("error creating theater"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateTheater(tx, theater)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error creating theater", err.Error())
	})
}
