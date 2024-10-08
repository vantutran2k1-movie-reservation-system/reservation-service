package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestGenreRepository_GetGenre(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	repo := NewGenreRepository(db)

	genre := utils.GenerateRandomGenre()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(genre.ID, genre.Name)

		mock.ExpectQuery(`SELECT \* FROM "genres" WHERE "genres"\."id" = \$1 ORDER BY "genres"\."id" LIMIT \$2`).
			WithArgs(genre.ID, 1).
			WillReturnRows(rows)

		result, err := repo.GetGenre(genre.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genre.ID, result.ID)
		assert.Equal(t, genre.Name, result.Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "genres" WHERE "genres"\."id" = \$1 ORDER BY "genres"\."id" LIMIT \$2`).
			WithArgs(genre.ID, 1).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetGenre(genre.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGenreRepository_CreateGenre(t *testing.T) {
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
