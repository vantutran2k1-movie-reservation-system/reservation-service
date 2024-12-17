package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
)

func TestTheaterLocationRepository_GetLocation(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterLocationRepository(db)

	location := utils.GenerateTheaterLocation()
	filter := filters.TheaterLocationFilter{
		Filter:    &filters.SingleFilter{},
		TheaterID: &filters.Condition{Operator: filters.OpEqual, Value: location.TheaterID},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theater_locations" WHERE theater_id = $1 ORDER BY "theater_locations"."id" LIMIT $2`)).
			WithArgs(filter.TheaterID.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(location))

		result, err := repo.GetLocation(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, location, result)
	})

	t.Run("location not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theater_locations" WHERE theater_id = $1 ORDER BY "theater_locations"."id" LIMIT $2`)).
			WithArgs(filter.TheaterID.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetLocation(filter)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting location", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theater_locations" WHERE theater_id = $1 ORDER BY "theater_locations"."id" LIMIT $2`)).
			WithArgs(filter.TheaterID.Value, 1).
			WillReturnError(errors.New("error getting location"))

		result, err := repo.GetLocation(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting location", err.Error())
	})
}

func TestTheaterLocationRepository_CreateTheaterLocation(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterLocationRepository(db)

	location := utils.GenerateTheaterLocation()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "theater_locations" ("id","theater_id","city_id","address","postal_code","latitude","longitude") VALUES ($1,$2,$3,$4,$5,$6,$7)`)).
			WithArgs(location.ID, location.TheaterID, location.CityID, location.Address, location.PostalCode, location.Latitude, location.Longitude).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateTheaterLocation(tx, location)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error creating location", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "theater_locations" ("id","theater_id","city_id","address","postal_code","latitude","longitude") VALUES ($1,$2,$3,$4,$5,$6,$7)`)).
			WithArgs(location.ID, location.TheaterID, location.CityID, location.Address, location.PostalCode, location.Latitude, location.Longitude).
			WillReturnError(errors.New("error creating location"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateTheaterLocation(tx, location)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error creating location", err.Error())
	})
}

func TestTheaterLocationRepository_UpdateTheaterLocation(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterLocationRepository(db)

	theater := utils.GenerateTheater()
	location := utils.GenerateTheaterLocation()
	location.TheaterID = &theater.ID

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "theater_locations" SET "theater_id"=$1,"city_id"=$2,"address"=$3,"postal_code"=$4,"latitude"=$5,"longitude"=$6 WHERE "id" = $7`)).
			WithArgs(location.TheaterID, location.CityID, location.Address, location.PostalCode, location.Latitude, location.Longitude, location.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.UpdateTheaterLocation(tx, location)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error updating location", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "theater_locations" SET "theater_id"=$1,"city_id"=$2,"address"=$3,"postal_code"=$4,"latitude"=$5,"longitude"=$6 WHERE "id" = $7`)).
			WithArgs(location.TheaterID, location.CityID, location.Address, location.PostalCode, location.Latitude, location.Longitude, location.ID).
			WillReturnError(errors.New("error updating location"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UpdateTheaterLocation(tx, location)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error updating location", err.Error())
	})
}
