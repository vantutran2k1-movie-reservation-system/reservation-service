package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestMovieRepository_CreateMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	movie := utils.GenerateSampleMovie()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "movies"`).
			WithArgs(movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, movie.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := NewMovieRepository(db).CreateMovie(tx, movie)
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "movies"`).
			WithArgs(movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, movie.UpdatedAt).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := NewMovieRepository(db).CreateMovie(tx, movie)
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
