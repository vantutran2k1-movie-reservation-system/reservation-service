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

func TestTheaterRepository_GetTheaterByName(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewTheaterRepository(db)

	theater := utils.GenerateTheater()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" WHERE name = $1 ORDER BY "theaters"."id" LIMIT $2`)).
			WithArgs(theater.Name, 1).
			WillReturnRows(utils.GenerateSqlMockRow(theater))

		result, err := repo.GetTheaterByName(theater.Name)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, theater, result)
	})

	t.Run("theater not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" WHERE name = $1 ORDER BY "theaters"."id" LIMIT $2`)).
			WithArgs(theater.Name, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetTheaterByName(theater.Name)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting theater", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "theaters" WHERE name = $1 ORDER BY "theaters"."id" LIMIT $2`)).
			WithArgs(theater.Name, 1).
			WillReturnError(errors.New("error getting theater"))

		result, err := repo.GetTheaterByName(theater.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting theater", err.Error())
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
