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

func TestCountryRepository_GetCountry(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCountryRepository(db)

	country := utils.GenerateCountry()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE id = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.ID, 1).
			WillReturnRows(utils.GenerateSqlMockRow(country))

		result, err := repo.GetCountry(country.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, country, result)
	})

	t.Run("country not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE id = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.ID, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetCountry(country.ID)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting country", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE id = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.ID, 1).
			WillReturnError(errors.New("error getting country"))

		result, err := repo.GetCountry(country.ID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting country", err.Error())
	})
}

func TestCountryRepository_GetCountryByName(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCountryRepository(db)

	country := utils.GenerateCountry()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE name = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.Name, 1).
			WillReturnRows(utils.GenerateSqlMockRow(country))

		result, err := repo.GetCountryByName(country.Name)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, country, result)
	})

	t.Run("country not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE name = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.Name, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetCountryByName(country.Name)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting country", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE name = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.Name, 1).
			WillReturnError(errors.New("error getting country"))

		result, err := repo.GetCountryByName(country.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting country", err.Error())
	})
}

func TestCountryRepository_GetCountryByCode(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCountryRepository(db)

	country := utils.GenerateCountry()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE code = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.Code, 1).
			WillReturnRows(utils.GenerateSqlMockRow(country))

		result, err := repo.GetCountryByCode(country.Code)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, country, result)
	})

	t.Run("country not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE code = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.Code, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetCountryByCode(country.Code)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting country", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries" WHERE code = $1 ORDER BY "countries"."id" LIMIT $2`)).
			WithArgs(country.Code, 1).
			WillReturnError(errors.New("error getting country"))

		result, err := repo.GetCountryByCode(country.Code)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting country", err.Error())
	})
}

func TestCountryRepository_GetCountries(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCountryRepository(db)

	t.Run("success", func(t *testing.T) {
		countries := utils.GenerateCountries(3)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries"`)).
			WillReturnRows(utils.GenerateSqlMockRows(countries))

		result, err := repo.GetCountries()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, countries, result)
	})

	t.Run("error getting countries", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "countries"`)).
			WillReturnError(errors.New("error getting countries"))

		result, err := repo.GetCountries()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting countries", err.Error())
	})
}

func TestCountryRepository_CreateCountry(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewCountryRepository(db)

	country := utils.GenerateCountry()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "countries" ("id","name","code") VALUES ($1,$2,$3)`)).
			WithArgs(country.ID, country.Name, country.Code).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateCountry(tx, country)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error creating country", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "countries" ("id","name","code") VALUES ($1,$2,$3)`)).
			WithArgs(country.ID, country.Name, country.Code).
			WillReturnError(errors.New("error creating country"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateCountry(tx, country)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error creating country", err.Error())
	})
}
