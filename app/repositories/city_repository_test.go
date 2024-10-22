package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
)

func TestCityRepository_GetCityByName(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCityRepository(db)

	city := utils.GenerateCity()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE state_id = $1 AND name = $2 ORDER BY "cities"."id" LIMIT $3`)).
			WithArgs(city.StateID, city.Name, 1).
			WillReturnRows(utils.GenerateSqlMockRow(city))

		result, err := repo.GetCityByName(city.StateID, city.Name)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, city, result)
	})

	t.Run("city not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE state_id = $1 AND name = $2 ORDER BY "cities"."id" LIMIT $3`)).
			WithArgs(city.StateID, city.Name, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetCityByName(city.StateID, city.Name)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting city", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE state_id = $1 AND name = $2 ORDER BY "cities"."id" LIMIT $3`)).
			WithArgs(city.StateID, city.Name, 1).
			WillReturnError(errors.New("error getting city"))

		result, err := repo.GetCityByName(city.StateID, city.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting city", err.Error())
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
