package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"gorm.io/gorm"
)

func TestMovieRepository_GetMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	repo := NewMovieRepository(db)

	expectedMovie := utils.GenerateRandomMovie()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "duration_minutes", "language", "rating", "created_at", "updated_at", "created_by", "last_updated_by"}).
			AddRow(expectedMovie.ID, expectedMovie.Title, expectedMovie.Description, expectedMovie.ReleaseDate, expectedMovie.DurationMinutes, expectedMovie.Language, expectedMovie.Rating, expectedMovie.CreatedAt, expectedMovie.UpdatedAt, expectedMovie.CreatedBy, expectedMovie.LastUpdatedBy)

		mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."id" = \$1 ORDER BY "movies"."id" LIMIT \$2`).
			WithArgs(expectedMovie.ID, 1).
			WillReturnRows(rows)

		movie, err := repo.GetMovie(expectedMovie.ID)

		assert.NotNil(t, movie)
		assert.Nil(t, err)

		assert.Equal(t, expectedMovie.ID, movie.ID)
		assert.Equal(t, expectedMovie.Title, movie.Title)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("movie not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."id" = \$1 ORDER BY "movies"."id" LIMIT \$2`).
			WithArgs(expectedMovie.ID, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		movie, err := repo.GetMovie(expectedMovie.ID)

		assert.Nil(t, movie)
		assert.NotNil(t, err)

		assert.Equal(t, gorm.ErrRecordNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "movies" WHERE "movies"\."id" = \$1 ORDER BY "movies"."id" LIMIT \$2`).
			WithArgs(expectedMovie.ID, 1).
			WillReturnError(errors.New("query error"))

		movie, err := repo.GetMovie(expectedMovie.ID)

		assert.Nil(t, movie)
		assert.Error(t, err)

		assert.Equal(t, "query error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMovieRepository_GetMovies(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	repo := NewMovieRepository(db)

	movies := make([]*models.Movie, 3)
	for i := 0; i < len(movies); i++ {
		movies[i] = utils.GenerateRandomMovie()
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "duration_minutes", "language", "rating", "created_at", "updated_at", "created_by", "last_updated_by"}).
			AddRow(movies[0].ID, movies[0].Title, movies[0].Description, movies[0].ReleaseDate, movies[0].DurationMinutes, movies[0].Language, movies[0].Rating, movies[0].CreatedAt, movies[0].UpdatedAt, movies[0].CreatedBy, movies[0].LastUpdatedBy).
			AddRow(movies[1].ID, movies[1].Title, movies[1].Description, movies[1].ReleaseDate, movies[1].DurationMinutes, movies[1].Language, movies[1].Rating, movies[1].CreatedAt, movies[1].UpdatedAt, movies[1].CreatedBy, movies[1].LastUpdatedBy)

		mock.ExpectQuery(`SELECT \* FROM "movies" LIMIT \$1 OFFSET \$2`).
			WithArgs(2, 2).
			WillReturnRows(rows)

		result, err := repo.GetMovies(2, 2)

		assert.NotNil(t, result)
		assert.Nil(t, err)

		assert.Equal(t, movies[0], result[0])
		assert.Equal(t, movies[1], result[1])

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "movies" LIMIT \$1 OFFSET \$2`).
			WithArgs(2, 2).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetMovies(2, 2)

		assert.Nil(t, result)
		assert.NotNil(t, err)

		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMovieRepository_GetNumbersOfMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	repo := NewMovieRepository(db)

	movies := make([]*models.Movie, 3)
	for i := 0; i < len(movies); i++ {
		movies[i] = utils.GenerateRandomMovie()
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "movies"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		result, err := repo.GetNumbersOfMovie()

		assert.NotNil(t, result)
		assert.Nil(t, err)

		assert.Equal(t, 2, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT count\(\*\) FROM "movies"`).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetNumbersOfMovie()

		assert.Equal(t, result, 0)
		assert.NotNil(t, err)

		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMovieRepository_CreateMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	movie := utils.GenerateRandomMovie()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "movies"`).
			WithArgs(movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, movie.UpdatedAt, movie.CreatedBy, movie.LastUpdatedBy).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := NewMovieRepository(db).CreateMovie(tx, movie)
		tx.Commit()

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "movies"`).
			WithArgs(movie.ID, movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, movie.UpdatedAt, movie.CreatedBy, movie.LastUpdatedBy).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := NewMovieRepository(db).CreateMovie(tx, movie)
		tx.Rollback()

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMovieRepository_UpdateMovie(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	repo := NewMovieRepository(db)

	movie := utils.GenerateRandomMovie()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "movies" SET "title"=\$1,"description"=\$2,"release_date"=\$3,"duration_minutes"=\$4,"language"=\$5,"rating"=\$6,"created_at"=\$7,"updated_at"=\$8,"created_by"=\$9,"last_updated_by"=\$10 WHERE "id" = \$11`).
			WithArgs(movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, sqlmock.AnyArg(), movie.CreatedBy, movie.LastUpdatedBy, movie.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.UpdateMovie(tx, movie)
		tx.Commit()

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "movies" SET "title"=\$1,"description"=\$2,"release_date"=\$3,"duration_minutes"=\$4,"language"=\$5,"rating"=\$6,"created_at"=\$7,"updated_at"=\$8,"created_by"=\$9,"last_updated_by"=\$10 WHERE "id" = \$11`).
			WithArgs(movie.Title, movie.Description, movie.ReleaseDate, movie.DurationMinutes, movie.Language, movie.Rating, movie.CreatedAt, sqlmock.AnyArg(), movie.CreatedBy, movie.LastUpdatedBy, movie.ID).
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UpdateMovie(tx, movie)
		tx.Rollback()

		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
