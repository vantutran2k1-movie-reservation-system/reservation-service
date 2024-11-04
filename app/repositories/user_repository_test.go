package repositories

import (
	"errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestUserRepository_GetUser(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserRepository(db)

	user := utils.GenerateUser()
	filter := filters.UserFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: user.ID},
		Email:  &filters.Condition{Operator: filters.OpEqual, Value: user.Email},
	}

	t.Run("success", func(t *testing.T) {
		rows := utils.GenerateSqlMockRow(user)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND email = $2 ORDER BY "users"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Email.Value, 1).
			WillReturnRows(rows)

		result, err := repo.GetUser(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND email = $2 ORDER BY "users"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Email.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetUser(filter)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting user", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND email = $2 ORDER BY "users"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Email.Value, 1).
			WillReturnError(errors.New("error getting user"))

		result, err := repo.GetUser(filter)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "error getting user", err.Error())
	})
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	user := utils.GenerateUser()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
			WithArgs(user.ID, user.Email, user.PasswordHash, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := NewUserRepository(db).CreateUser(tx, user)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash","created_at","updated_at") VALUES ($1,$2,$3,$4,$5)`)).
			WithArgs(user.ID, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := NewUserRepository(db).CreateUser(tx, user)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	user := utils.GenerateUser()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "password_hash"=$1,"updated_at"=$2 WHERE "id" = $3`)).
			WithArgs(user.PasswordHash, sqlmock.AnyArg(), user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		result, err := NewUserRepository(db).UpdatePassword(tx, user, user.PasswordHash)
		tx.Commit()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.PasswordHash, result.PasswordHash)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "password_hash"=$1,"updated_at"=$2 WHERE "id" = $3`)).
			WithArgs(user.PasswordHash, sqlmock.AnyArg(), user.ID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		result, err := NewUserRepository(db).UpdatePassword(tx, user, user.PasswordHash)
		tx.Rollback()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
