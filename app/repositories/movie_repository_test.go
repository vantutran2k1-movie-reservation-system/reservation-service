package repositories

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestMovieRepository_GetMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewMovieRepository(db)

	movie := utils.GenerateRandomMovie()
	genres := make([]*models.Genre, 3)
	for i := 0; i < len(genres); i++ {
		genres[i] = utils.GenerateRandomGenre()
	}

	t.Run("success with genres", func(t *testing.T) {
		movieRows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "duration_minutes", "language", "rating", "created_at", "updated_at", "created_by", "last_updated_by"}).
			AddRow(movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, movie.UpdatedAt, movie.CreatedBy, movie.LastUpdatedBy)
		genreRows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(genres[0].ID, genres[0].Name).
			AddRow(genres[1].ID, genres[1].Name).
			AddRow(genres[2].ID, genres[2].Name)
		movieGenreRows := sqlmock.NewRows([]string{"movie_id", "genre_id"}).
			AddRow(movie.ID, genres[0].ID).
			AddRow(movie.ID, genres[1].ID).
			AddRow(movie.ID, genres[2].ID)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" WHERE id = $1 ORDER BY "movies"."id" LIMIT $2`)).
			WithArgs(movie.ID, 1).
			WillReturnRows(movieRows)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movie_genres" WHERE "movie_genres"."movie_id" = $1`)).
			WithArgs(movie.ID).
			WillReturnRows(movieGenreRows)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "genres" WHERE "genres"."id" IN ($1,$2,$3)`)).
			WithArgs(genres[0].ID, genres[1].ID, genres[2].ID).
			WillReturnRows(genreRows)

		result, err := repo.GetMovie(movie.ID, true)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movie.ID, result.ID)
		assert.Equal(t, movie.Title, result.Title)
		assert.Equal(t, genres[0], &result.Genres[0])
		assert.Equal(t, genres[1], &result.Genres[1])
		assert.Equal(t, genres[2], &result.Genres[2])
	})

	t.Run("success without genres", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "duration_minutes", "language", "rating", "created_at", "updated_at", "created_by", "last_updated_by"}).
			AddRow(movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, movie.UpdatedAt, movie.CreatedBy, movie.LastUpdatedBy)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" WHERE id = $1 ORDER BY "movies"."id" LIMIT $2`)).
			WithArgs(movie.ID, 1).
			WillReturnRows(rows)

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

	movies := make([]*models.Movie, 3)
	for i := 0; i < len(movies); i++ {
		movies[i] = utils.GenerateRandomMovie()
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "duration_minutes", "language", "rating", "created_at", "updated_at", "created_by", "last_updated_by"}).
			AddRow(movies[0].ID, movies[0].Title, movies[0].Description, movies[0].ReleaseDate, movies[0].DurationMinutes, movies[0].Language, movies[0].Rating, movies[0].CreatedAt, movies[0].UpdatedAt, movies[0].CreatedBy, movies[0].LastUpdatedBy).
			AddRow(movies[1].ID, movies[1].Title, movies[1].Description, movies[1].ReleaseDate, movies[1].DurationMinutes, movies[1].Language, movies[1].Rating, movies[1].CreatedAt, movies[1].UpdatedAt, movies[1].CreatedBy, movies[1].LastUpdatedBy).
			AddRow(movies[2].ID, movies[2].Title, movies[2].Description, movies[2].ReleaseDate, movies[2].DurationMinutes, movies[2].Language, movies[2].Rating, movies[2].CreatedAt, movies[2].UpdatedAt, movies[2].CreatedBy, movies[2].LastUpdatedBy)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "movies" LIMIT $1 OFFSET $2`)).
			WithArgs(2, 2).
			WillReturnRows(rows)

		result, err := repo.GetMovies(2, 2)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, movies[0], result[0])
		assert.Equal(t, movies[1], result[1])
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

	movies := make([]*models.Movie, 3)
	for i := 0; i < len(movies); i++ {
		movies[i] = utils.GenerateRandomMovie()
	}

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

	movie := utils.GenerateRandomMovie()

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

	movie := utils.GenerateRandomMovie()

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
