package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestGenreRepository_GetGenre(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewGenreRepository(db)

	genre := utils.GenerateGenre()
	filter := filters.GenreFilter{
		Filter: &filters.SingleFilter{Logic: filters.And},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: genre.ID.String()},
		Name:   &filters.Condition{Operator: filters.OpEqual, Value: genre.Name},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "genres" WHERE id = $1 AND name = $2 ORDER BY "genres"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Name.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(genre))

		result, err := repo.GetGenre(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genre, result)
	})

	t.Run("error getting genre", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "genres" WHERE id = $1 AND name = $2 ORDER BY "genres"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Name.Value, 1).
			WillReturnError(errors.New("error getting genre"))

		result, err := repo.GetGenre(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting genre", err.Error())
	})
}

func TestGenreRepository_GetGenres(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewGenreRepository(db)

	limit := 10
	offset := 10
	filter := filters.GenreFilter{Filter: &filters.MultiFilter{Logic: filters.And, Limit: &limit, Offset: &offset}}

	t.Run("success", func(t *testing.T) {
		genres := make([]*models.Genre, 3)
		for i := 0; i < len(genres); i++ {
			genres[i] = utils.GenerateGenre()
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "genres" LIMIT $1 OFFSET $2`)).
			WithArgs(limit, offset).
			WillReturnRows(utils.GenerateSqlMockRows(genres))

		result, err := repo.GetGenres(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, genres, result)
	})

	t.Run("error getting genres", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "genres"`)).
			WillReturnError(errors.New("error getting genres"))

		result, err := repo.GetGenres(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting genres", err.Error())
	})
}

func TestGenreRepository_GetGenreIDs(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewGenreRepository(db)

	limit := 10
	offset := 10
	filter := filters.GenreFilter{Filter: &filters.MultiFilter{Logic: filters.And, Limit: &limit, Offset: &offset}}

	t.Run("success", func(t *testing.T) {
		genres := make([]*models.Genre, 3)
		for i := 0; i < len(genres); i++ {
			genres[i] = utils.GenerateGenre()
		}

		rows := sqlmock.NewRows([]string{"id"})
		for _, genre := range genres {
			rows.AddRow(genre.ID)
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "genres" LIMIT $1 OFFSET $2`)).
			WithArgs(limit, offset).
			WillReturnRows(rows)

		result, err := repo.GetGenreIDs(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, len(genres), len(result))
		for i, g := range genres {
			assert.Equal(t, g.ID, result[i])
		}
	})

	t.Run("error getting genres", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "genres" LIMIT $1 OFFSET $2`)).
			WithArgs(limit, offset).
			WillReturnError(errors.New("error getting genres"))

		result, err := repo.GetGenreIDs(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting genres", err.Error())
	})
}

func TestGenreRepository_CreateGenre(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewGenreRepository(db)

	genre := utils.GenerateGenre()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "genres" ("id","name") VALUES ($1,$2)`)).
			WithArgs(genre.ID, genre.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateGenre(tx, genre)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error creating genre", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "genres" ("id","name") VALUES ($1,$2)`)).
			WithArgs(genre.ID, genre.Name).WillReturnError(errors.New("error creating genre"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateGenre(tx, genre)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error creating genre", err.Error())
	})
}

func TestGenreRepository_UpdateGenre(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewGenreRepository(db)

	genre := utils.GenerateGenre()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "genres" SET "name"=$1 WHERE "id" = $2`)).
			WithArgs(genre.Name, genre.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.UpdateGenre(tx, genre)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error updating genre", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "genres" SET "name"=$1 WHERE "id" = $2`)).
			WithArgs(genre.Name, genre.ID).
			WillReturnError(errors.New("error updating genre"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UpdateGenre(tx, genre)
		tx.Rollback()

		assert.NotNil(t, err)
	})
}

func TestGenreRepository_DeleteGenre(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewGenreRepository(db)

	genre := utils.GenerateGenre()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "genres" WHERE "genres"."id" = $1`)).
			WithArgs(genre.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.DeleteGenre(tx, genre)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error deleting genre", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "genres" WHERE "genres"."id" = $1`)).
			WithArgs(genre.ID).
			WillReturnError(errors.New("error deleting genre"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.DeleteGenre(tx, genre)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error deleting genre", err.Error())
	})
}
