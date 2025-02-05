package repositories

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
	"time"
)

func TestShowRepository_GetShow(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewShowRepository(db)

	show := utils.GenerateShow()
	filter := filters.ShowFilter{
		Filter: &filters.SingleFilter{},
		Id:     &filters.Condition{Operator: filters.OpEqual, Value: show.Id},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "shows" WHERE id = $1 ORDER BY "shows"."id" LIMIT $2`)).
			WithArgs(show.Id, 1).
			WillReturnRows(utils.GenerateSqlMockRow(show))

		result, err := repo.GetShow(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	})

	t.Run("show not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "shows" WHERE id = $1 ORDER BY "shows"."id" LIMIT $2`)).
			WithArgs(show.Id, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetShow(filter)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting show", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "shows" WHERE id = $1 ORDER BY "shows"."id" LIMIT $2`)).
			WithArgs(show.Id, 1).
			WillReturnError(errors.New("error getting show"))

		result, err := repo.GetShow(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting show")
	})
}

func TestShowRepository_GetShows(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewShowRepository(db)

	shows := utils.GenerateShows(3)
	filter := filters.ShowFilter{
		Filter: &filters.SingleFilter{},
		Status: &filters.Condition{Operator: filters.OpEqual, Value: constants.Active},
	}

	query := regexp.QuoteMeta(`SELECT * FROM "shows" WHERE status = $1`)
	args := []driver.Value{constants.Active}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(args...).
			WillReturnRows(utils.GenerateSqlMockRows(shows))

		result, err := repo.GetShows(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, shows, result)
	})

	t.Run("error getting shows", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(args...).
			WillReturnError(errors.New("error getting shows"))

		result, err := repo.GetShows(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting shows")
	})
}

func TestShowRepository_IsShowInValidTimeRange(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewShowRepository(db)

	show := utils.GenerateShow()
	query := regexp.QuoteMeta(`
		SELECT id 
		FROM shows
		WHERE TRUE
		  	AND status IN ($1, $2)
		  	AND theater_id = $3 
		  	AND (
		  	    start_time BETWEEN $4 AND $5 
		  	    OR end_time BETWEEN $6 AND $7
		  	    OR (start_time <= $8 AND end_time >= $9)
		  	)
`)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(constants.Active, constants.Scheduled, show.TheaterId, show.StartTime, show.EndTime, show.StartTime, show.EndTime, show.StartTime, show.EndTime).
			WillReturnRows(utils.GenerateSqlMockRow(nil))

		result, err := repo.IsShowInValidTimeRange(*show.TheaterId, show.StartTime, show.EndTime)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("not valid", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(constants.Active, constants.Scheduled, show.TheaterId, show.StartTime, show.EndTime, show.StartTime, show.EndTime, show.StartTime, show.EndTime).
			WillReturnRows(utils.GenerateSqlMockRow(show))

		result, err := repo.IsShowInValidTimeRange(*show.TheaterId, show.StartTime, show.EndTime)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("error getting show", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(constants.Active, constants.Scheduled, show.TheaterId, show.StartTime, show.EndTime, show.StartTime, show.EndTime, show.StartTime, show.EndTime).
			WillReturnError(errors.New("error getting show"))

		result, err := repo.IsShowInValidTimeRange(*show.TheaterId, show.StartTime, show.EndTime)

		assert.NotNil(t, result)
		assert.NotNil(t, err)
		assert.False(t, result)
		assert.EqualError(t, err, "error getting show")
	})
}

func TestShowRepository_CreateShow(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewShowRepository(db)

	show := utils.GenerateShow()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "shows" ("id","movie_id","theater_id","start_time","end_time","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).
			WithArgs(show.Id, show.MovieId, show.TheaterId, show.StartTime, show.EndTime, show.Status, show.CreatedAt, show.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateShow(tx, show)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "shows" ("id","movie_id","theater_id","start_time","end_time","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).
			WithArgs(show.Id, show.MovieId, show.TheaterId, show.StartTime, show.EndTime, show.Status, show.CreatedAt, show.UpdatedAt).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateShow(tx, show)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "db error")
	})
}

func TestShowRepository_UpdateShowStatus(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewShowRepository(db)

	show := utils.GenerateShow()
	show.Status = constants.Scheduled

	statement := regexp.QuoteMeta(`UPDATE "shows" SET "status"=$1,"updated_at"=$2 WHERE id = $3`)
	args := []driver.Value{constants.Active, sqlmock.AnyArg(), show.Id}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(statement).WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.UpdateShowStatus(tx, show.Id, constants.Active)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error update show", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(statement).WithArgs(args...).WillReturnError(errors.New("error updating show"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UpdateShowStatus(tx, show.Id, constants.Active)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error updating show")
	})
}

func TestShowRepository_ScheduleActivateShows(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewShowRepository(db)

	statement := regexp.QuoteMeta(`UPDATE "shows" SET "status"=$1,"updated_at"=$2 WHERE start_time <= $3 AND status = $4`)
	args := []driver.Value{constants.Active, sqlmock.AnyArg(), sqlmock.AnyArg(), constants.Scheduled}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(statement).WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.ScheduleActivateShows(tx, time.Hour*24)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error updating shows", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(statement).WithArgs(args...).WillReturnError(errors.New("error updating shows"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.ScheduleActivateShows(tx, time.Hour*24)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error updating shows")
	})
}

func TestShowRepository_ScheduleCompleteShows(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewShowRepository(db)

	statement := regexp.QuoteMeta(`UPDATE "shows" SET "status"=$1,"updated_at"=$2 WHERE end_time <= $3 AND status = $4`)
	args := []driver.Value{constants.Completed, sqlmock.AnyArg(), sqlmock.AnyArg(), constants.Active}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(statement).WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.ScheduleCompleteShows(tx)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error updating shows", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(statement).WithArgs(args...).WillReturnError(errors.New("error updating shows"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.ScheduleCompleteShows(tx)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error updating shows")
	})
}
