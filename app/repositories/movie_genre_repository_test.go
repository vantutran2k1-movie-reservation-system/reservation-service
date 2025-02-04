package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
)

func TestMovieGenreRepository_UpdateGenresOfMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieGenreRepository(db)

	movie := utils.GenerateMovie()
	genreIDs := make([]uuid.UUID, 3)
	for i := 0; i < len(genreIDs); i++ {
		genreIDs[i] = utils.GenerateGenre().ID
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "movie_genres" WHERE movie_id = $1`)).
			WithArgs(movie.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "movie_genres" ("movie_id","genre_id") VALUES ($1,$2),($3,$4),($5,$6)`)).
			WithArgs(movie.ID, genreIDs[0], movie.ID, genreIDs[1], movie.ID, genreIDs[2]).
			WillReturnResult(sqlmock.NewResult(3, 3))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.UpdateGenresOfMovie(tx, movie.ID, genreIDs)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("no genres", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "movie_genres" WHERE movie_id = $1`)).
			WithArgs(movie.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		var emptyGenres []uuid.UUID
		err := repo.UpdateGenresOfMovie(tx, movie.ID, emptyGenres)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error deleting genres", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "movie_genres" WHERE movie_id = $1`)).
			WithArgs(movie.ID).
			WillReturnError(errors.New("error deleting genres"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UpdateGenresOfMovie(tx, movie.ID, genreIDs)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error deleting genres", err.Error())
	})

	t.Run("error inserting genres", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "movie_genres" WHERE movie_id = $1`)).
			WithArgs(movie.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "movie_genres" ("movie_id","genre_id") VALUES ($1,$2),($3,$4),($5,$6)`)).
			WithArgs(movie.ID, genreIDs[0], movie.ID, genreIDs[1], movie.ID, genreIDs[2]).
			WillReturnError(errors.New("error inserting genres"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UpdateGenresOfMovie(tx, movie.ID, genreIDs)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error inserting genres", err.Error())
	})
}

func TestMovieGenreRepository_DeleteByMovieId(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieGenreRepository(db)

	movie := utils.GenerateMovie()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "movie_genres" WHERE movie_id = $1`)).
			WithArgs(movie.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.DeleteByMovieId(tx, movie.ID)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error deleting by movie id", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "movie_genres" WHERE movie_id = $1`)).
			WithArgs(movie.ID).
			WillReturnError(errors.New("error deleting by movie id"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.DeleteByMovieId(tx, movie.ID)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error deleting by movie id", err.Error())
	})
}

func TestMovieGenreRepository_DeleteByGenreId(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieGenreRepository(db)

	genre := utils.GenerateGenre()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "movie_genres" WHERE genre_id = $1`)).
			WithArgs(genre.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.DeleteByGenreId(tx, genre.ID)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error deleting by genre id", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "movie_genres" WHERE genre_id = $1`)).
			WithArgs(genre.ID).
			WillReturnError(errors.New("error deleting by genre id"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.DeleteByGenreId(tx, genre.ID)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error deleting by genre id", err.Error())
	})
}
