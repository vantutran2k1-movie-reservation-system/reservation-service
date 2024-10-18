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

func TestStateRepository_GetStateByName(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewStateRepository(db)

	state := utils.GenerateState()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "code", "country_id"}).
			AddRow(state.ID, state.Name, state.Code, state.CountryID)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE country_id = $1 AND name = $2 ORDER BY "states"."id" LIMIT $3`)).
			WithArgs(state.CountryID, state.Name, 1).
			WillReturnRows(rows)

		result, err := repo.GetStateByName(state.CountryID, state.Name)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, state, result)
	})

	t.Run("state not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE country_id = $1 AND name = $2 ORDER BY "states"."id" LIMIT $3`)).
			WithArgs(state.CountryID, state.Name, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetStateByName(state.CountryID, state.Name)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting state", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE country_id = $1 AND name = $2 ORDER BY "states"."id" LIMIT $3`)).
			WithArgs(state.CountryID, state.Name, 1).
			WillReturnError(errors.New("error getting state"))

		result, err := repo.GetStateByName(state.CountryID, state.Name)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting state", err.Error())
	})
}

func TestStateRepository_GetStatesByCountry(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewStateRepository(db)

	countryID := uuid.New()

	t.Run("success", func(t *testing.T) {
		states := utils.GenerateStates(3)

		rows := sqlmock.NewRows([]string{"id", "name", "code", "country_id"})
		for _, state := range states {
			rows.AddRow(state.ID, state.Name, state.Code, state.CountryID)
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE country_id = $1`)).
			WithArgs(countryID).
			WillReturnRows(rows)

		result, err := repo.GetStatesByCountry(countryID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, states, result)
	})

	t.Run("error getting states", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE country_id = $1`)).
			WithArgs(countryID).
			WillReturnError(errors.New("error getting states"))

		result, err := repo.GetStatesByCountry(countryID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting states", err.Error())
	})
}

func TestStateRepository_CreateState(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewStateRepository(db)

	state := utils.GenerateState()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "states" ("id","name","code","country_id") VALUES ($1,$2,$3,$4)`)).
			WithArgs(state.ID, state.Name, state.Code, state.CountryID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateState(tx, state)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error creating state", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "states" ("id","name","code","country_id") VALUES ($1,$2,$3,$4)`)).
			WithArgs(state.ID, state.Name, state.Code, state.CountryID).
			WillReturnError(errors.New("error creating state"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateState(tx, state)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "error creating state", err.Error())
	})
}
