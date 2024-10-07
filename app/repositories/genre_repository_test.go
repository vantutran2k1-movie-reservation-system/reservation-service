package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestGenreRepository_CreateUser(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	genre := utils.GenerateRandomGenre()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "genres"`).WithArgs(genre.ID, genre.Name).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := NewGenreRepository(db).CreateGenre(tx, genre)
		tx.Commit()

		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "genres"`).WithArgs(genre.ID, genre.Name).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := NewGenreRepository(db).CreateGenre(tx, genre)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
