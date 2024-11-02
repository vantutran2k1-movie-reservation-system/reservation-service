package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
)

func TestCityRepository_GetCity(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCityRepository(db)

	city := utils.GenerateCity()
	filter := filters.CityFilter{
		Filter:  &filters.SingleFilter{Logic: filters.And},
		ID:      &filters.Condition{Operator: filters.OpEqual, Value: city.ID},
		StateID: &filters.Condition{Operator: filters.OpEqual, Value: city.StateID},
		Name:    &filters.Condition{Operator: filters.OpEqual, Value: city.Name},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1 AND state_id = $2 AND name = $3 ORDER BY "cities"."id" LIMIT $4`)).
			WithArgs(filter.ID.Value, filter.StateID.Value, filter.Name.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(city))

		result, err := repo.GetCity(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, city, result)
	})

	t.Run("city not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1 AND state_id = $2 AND name = $3 ORDER BY "cities"."id" LIMIT $4`)).
			WithArgs(filter.ID.Value, filter.StateID.Value, filter.Name.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetCity(filter)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting city", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1 AND state_id = $2 AND name = $3 ORDER BY "cities"."id" LIMIT $4`)).
			WithArgs(filter.ID.Value, filter.StateID.Value, filter.Name.Value, 1).
			WillReturnError(errors.New("error getting city"))

		result, err := repo.GetCity(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting city", err.Error())
	})
}

func TestCityRepository_GetCities(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCityRepository(db)

	cities := utils.GenerateCities(3)
	filter := filters.CityFilter{
		Filter:  &filters.MultiFilter{Logic: filters.And},
		StateID: &filters.Condition{Operator: filters.OpEqual, Value: uuid.New()},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE state_id = $1`)).
			WithArgs(filter.StateID.Value).
			WillReturnRows(utils.GenerateSqlMockRows(cities))

		result, err := repo.GetCities(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, cities, result)
	})

	t.Run("error getting cities", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE state_id = $1`)).
			WithArgs(filter.StateID.Value).
			WillReturnError(errors.New("error getting cities"))

		result, err := repo.GetCities(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting cities", err.Error())
	})
}

func TestCityRepository_CreateCity(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCityRepository(db)

	city := utils.GenerateCity()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "cities" ("id","name","state_id") VALUES ($1,$2,$3)`)).
			WithArgs(city.ID, city.Name, city.StateID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateCity(tx, city)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error creating city", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "cities" ("id","name","state_id") VALUES ($1,$2,$3)`)).
			WithArgs(city.ID, city.Name, city.StateID).
			WillReturnError(errors.New("error creating city"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateCity(tx, city)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error creating city", err.Error())
	})
}
