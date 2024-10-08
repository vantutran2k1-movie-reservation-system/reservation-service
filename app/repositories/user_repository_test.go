package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"gorm.io/gorm"
)

func TestUserRepository_GetUser(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	user := utils.GenerateRandomUser()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(user.ID, user.Email)

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 ORDER BY "users"."id" LIMIT \$2`).
			WithArgs(user.ID, 1).
			WillReturnRows(rows)

		result, err := NewUserRepository(db).GetUser(user.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(user.ID, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := NewUserRepository(db).GetUser(user.ID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(user.ID, 1).
			WillReturnError(errors.New("db error"))

		result, err := NewUserRepository(db).GetUser(user.ID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_FindUserByEmail(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	user := utils.GenerateRandomUser()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(user.ID, user.Email)

		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."email" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(user.Email, 1).
			WillReturnRows(rows)

		result, err := NewUserRepository(db).FindUserByEmail(user.Email)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."email" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(user.Email, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := NewUserRepository(db).FindUserByEmail(user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."email" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
			WithArgs(user.Email, 1).
			WillReturnError(errors.New("db error"))

		result, err := NewUserRepository(db).FindUserByEmail(user.Email)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	user := utils.GenerateRandomUser()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "users"`).
			WithArgs(user.ID, user.Email, user.PasswordHash, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := NewUserRepository(db).CreateUser(tx, user)
		tx.Commit()

		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "users"`).
			WithArgs(user.ID, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := NewUserRepository(db).CreateUser(tx, user)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	user := utils.GenerateRandomUser()
	hashedPassword := utils.GenerateRandomHashedPassword()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "password_hash"=\$1,"updated_at"=\$2 WHERE "id" = \$3`).
			WithArgs(hashedPassword, sqlmock.AnyArg(), user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		result, err := NewUserRepository(db).UpdatePassword(tx, user, hashedPassword)
		tx.Commit()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, hashedPassword, result.PasswordHash)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "password_hash"=\$1,"updated_at"=\$2 WHERE "id" = \$3`).
			WithArgs(user.PasswordHash, sqlmock.AnyArg(), user.ID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		result, err := NewUserRepository(db).UpdatePassword(tx, user, user.PasswordHash)
		tx.Rollback()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
