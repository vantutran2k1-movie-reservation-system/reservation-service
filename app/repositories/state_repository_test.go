package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
)

func TestStateRepository_GetState(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewStateRepository(db)

	state := utils.GenerateState()
	filter := filters.StateFilter{
		Filter:    &filters.SingleFilter{Logic: filters.And},
		ID:        &filters.Condition{Operator: filters.OpEqual, Value: state.ID},
		CountryID: &filters.Condition{Operator: filters.OpEqual, Value: state.CountryID},
		Name:      &filters.Condition{Operator: filters.OpEqual, Value: state.Name},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE id = $1 AND country_id = $2 AND name = $3 ORDER BY "states"."id" LIMIT $4`)).
			WithArgs(filter.ID.Value, filter.CountryID.Value, filter.Name.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(state))

		result, err := repo.GetState(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, state, result)
	})

	t.Run("state not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE id = $1 AND country_id = $2 AND name = $3 ORDER BY "states"."id" LIMIT $4`)).
			WithArgs(filter.ID.Value, filter.CountryID.Value, filter.Name.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetState(filter)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting state", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE id = $1 AND country_id = $2 AND name = $3 ORDER BY "states"."id" LIMIT $4`)).
			WithArgs(filter.ID.Value, filter.CountryID.Value, filter.Name.Value, 1).
			WillReturnError(errors.New("error getting state"))

		result, err := repo.GetState(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting state", err.Error())
	})
}

func TestStateRepository_GetStates(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewStateRepository(db)

	filter := filters.StateFilter{
		Filter:    &filters.MultiFilter{Logic: filters.And},
		CountryID: &filters.Condition{Operator: filters.OpEqual, Value: uuid.New()},
	}

	t.Run("success", func(t *testing.T) {
		states := utils.GenerateStates(3)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE country_id = $1`)).
			WithArgs(filter.CountryID.Value).
			WillReturnRows(utils.GenerateSqlMockRows(states))

		result, err := repo.GetStates(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, states, result)
	})

	t.Run("error getting states", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "states" WHERE country_id = $1`)).
			WithArgs(filter.CountryID.Value).
			WillReturnError(errors.New("error getting states"))

		result, err := repo.GetStates(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting states", err.Error())
	})
}

func TestStateRepository_CreateState(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
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
