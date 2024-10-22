package repositories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestMovieRepository_GetMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieRepository(db)

	movie := utils.GenerateMovie()
	genres := utils.GenerateGenres(3)

	t.Run("success with genres", func(t *testing.T) {
		movieGenreRows := sqlmock.NewRows([]string{"movie_id", "genre_id"})
		for _, genre := range genres {
			movieGenreRows.AddRow(movie.ID, genre.ID)
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" WHERE id = $1 ORDER BY "movies"."id" LIMIT $2`)).
			WithArgs(movie.ID, 1).
			WillReturnRows(utils.GenerateSqlMockRow(movie))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movie_genres" WHERE "movie_genres"."movie_id" = $1`)).
			WithArgs(movie.ID).
			WillReturnRows(movieGenreRows)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "genres" WHERE "genres"."id" IN ($1,$2,$3)`)).
			WithArgs(genres[0].ID, genres[1].ID, genres[2].ID).
			WillReturnRows(utils.GenerateSqlMockRows(genres))

		result, err := repo.GetMovie(movie.ID, true)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie.ID, result.ID)
		assert.Equal(t, movie.Title, result.Title)
		for i, genre := range genres {
			assert.Equal(t, genre, &result.Genres[i])
		}
	})

	t.Run("success without genres", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" WHERE id = $1 ORDER BY "movies"."id" LIMIT $2`)).
			WithArgs(movie.ID, 1).
			WillReturnRows(utils.GenerateSqlMockRow(movie))

		result, err := repo.GetMovie(movie.ID, false)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie, result)
	})

	t.Run("movie not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" WHERE id = $1 ORDER BY "movies"."id" LIMIT $2`)).
			WithArgs(movie.ID, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetMovie(movie.ID, false)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting movie", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" WHERE id = $1 ORDER BY "movies"."id" LIMIT $2`)).
			WithArgs(movie.ID, 1).
			WillReturnError(errors.New("error getting movie"))

		result, err := repo.GetMovie(movie.ID, false)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "error getting movie", err.Error())
	})
}

func TestMovieRepository_GetMovies(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieRepository(db)

	movies := utils.GenerateMovies(3)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" LIMIT $1 OFFSET $2`)).
			WithArgs(2, 2).
			WillReturnRows(utils.GenerateSqlMockRows(movies))

		result, err := repo.GetMovies(2, 2)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		for i, movie := range movies {
			assert.Equal(t, movie, result[i])
		}
	})

	t.Run("error getting movies", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" LIMIT $1 OFFSET $2`)).
			WithArgs(2, 2).
			WillReturnError(errors.New("error getting movies"))

		result, err := repo.GetMovies(2, 2)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting movies", err.Error())
	})
}

func TestMovieRepository_GetNumbersOfMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieRepository(db)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "movies"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		result, err := repo.GetNumbersOfMovie()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, 3, result)
	})

	t.Run("error counting movies", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "movies"`)).
			WillReturnError(errors.New("error counting movies"))

		result, err := repo.GetNumbersOfMovie()

		assert.Equal(t, result, 0)
		assert.NotNil(t, err)
		assert.Equal(t, "error counting movies", err.Error())
	})
}

func TestMovieRepository_CreateMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieRepository(db)

	movie := utils.GenerateMovie()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "movies" ("id","title","description","release_date","duration_minutes","language","rating","created_at","updated_at","created_by","last_updated_by") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)).
			WithArgs(movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, movie.UpdatedAt, movie.CreatedBy, movie.LastUpdatedBy).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateMovie(tx, movie)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error creating movie", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "movies" ("id","title","description","release_date","duration_minutes","language","rating","created_at","updated_at","created_by","last_updated_by") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)).
			WithArgs(movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, movie.UpdatedAt, movie.CreatedBy, movie.LastUpdatedBy).
			WillReturnError(errors.New("error creating movie"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateMovie(tx, movie)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error creating movie", err.Error())
	})
}

func TestMovieRepository_UpdateMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieRepository(db)

	movie := utils.GenerateMovie()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "movies" SET "title"=$1,"description"=$2,"release_date"=$3,"duration_minutes"=$4,"language"=$5,"rating"=$6,"created_at"=$7,"updated_at"=$8,"created_by"=$9,"last_updated_by"=$10 WHERE "id" = $11`)).
			WithArgs(movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, sqlmock.AnyArg(), movie.CreatedBy, movie.LastUpdatedBy, movie.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.UpdateMovie(tx, movie)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error updating movie", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "movies" SET "title"=$1,"description"=$2,"release_date"=$3,"duration_minutes"=$4,"language"=$5,"rating"=$6,"created_at"=$7,"updated_at"=$8,"created_by"=$9,"last_updated_by"=$10 WHERE "id" = $11`)).
			WithArgs(movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, sqlmock.AnyArg(), movie.CreatedBy, movie.LastUpdatedBy, movie.ID).
			WillReturnError(errors.New("error updating movie"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UpdateMovie(tx, movie)
		tx.Rollback()

		assert.Error(t, err)
		assert.Equal(t, "error updating movie", err.Error())
	})
}
