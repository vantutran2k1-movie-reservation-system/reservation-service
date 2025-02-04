package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
)

func TestSeatRepository_GetSeat(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewSeatRepository(db)

	seat := utils.GenerateSeat()
	filter := filters.SeatFilter{
		Filter:    &filters.SingleFilter{},
		TheaterId: &filters.Condition{Operator: filters.OpEqual, Value: seat.TheaterId},
		Row:       &filters.Condition{Operator: filters.OpEqual, Value: seat.Row},
		Number:    &filters.Condition{Operator: filters.OpEqual, Value: seat.Number},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "seats" WHERE theater_id = $1 AND row = $2 AND number = $3 ORDER BY "seats"."id" LIMIT $4`)).
			WithArgs(filter.TheaterId.Value, filter.Row.Value, filter.Number.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(seat))

		result, err := repo.GetSeat(filter)

		assert.NotNil(t, result)
		assert.NoError(t, err)
		assert.Equal(t, seat.TheaterId, result.TheaterId)
		assert.Equal(t, seat.Row, result.Row)
		assert.Equal(t, seat.Number, result.Number)
		assert.Equal(t, seat.Type, result.Type)
	})

	t.Run("seat not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "seats" WHERE theater_id = $1 AND row = $2 AND number = $3 ORDER BY "seats"."id" LIMIT $4`)).
			WithArgs(filter.TheaterId.Value, filter.Row.Value, filter.Number.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetSeat(filter)

		assert.Nil(t, result)
		assert.NoError(t, err)
	})

	t.Run("error getting seat", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "seats" WHERE theater_id = $1 AND row = $2 AND number = $3 ORDER BY "seats"."id" LIMIT $4`)).
			WithArgs(filter.TheaterId.Value, filter.Row.Value, filter.Number.Value, 1).
			WillReturnError(errors.New("error getting seat"))

		result, err := repo.GetSeat(filter)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "error getting seat")
	})
}

func TestSeatRepository_CreateSeat(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewSeatRepository(db)

	seat := utils.GenerateSeat()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "seats" ("id","theater_id","row","number","type") VALUES ($1,$2,$3,$4,$5)`)).
			WithArgs(seat.Id, seat.TheaterId, seat.Row, seat.Number, seat.Type).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateSeat(tx, seat)
		tx.Commit()

		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "seats" ("id","theater_id","row","number","type") VALUES ($1,$2,$3,$4,$5)`)).
			WithArgs(seat.Id, seat.TheaterId, seat.Row, seat.Number, seat.Type).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateSeat(tx, seat)
		tx.Rollback()

		assert.EqualError(t, err, "db error")
	})
}
