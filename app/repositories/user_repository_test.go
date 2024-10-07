package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"gorm.io/gorm"
)

func TestUserRepository_GetUser_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	expectedUser := utils.GenerateRandomUser()

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(expectedUser.ID, expectedUser.Email)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(expectedUser.ID, 1).
		WillReturnRows(rows)

	user, err := NewUserRepository(db).GetUser(expectedUser.ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUser_UserNotFound(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	user, err := NewUserRepository(db).GetUser(userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUser_QueryError(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnError(errors.New("query error"))

	user, err := NewUserRepository(db).GetUser(userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "query error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindUserByEmail_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	expectedUser := utils.GenerateRandomUser()

	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(expectedUser.ID, expectedUser.Email)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."email" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(expectedUser.Email, 1).
		WillReturnRows(rows)

	user, err := NewUserRepository(db).FindUserByEmail(expectedUser.Email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindUserByEmail_UserNotFound(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	email := "notfound@example.com"

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."email" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnRows(sqlmock.NewRows(nil))

	user, err := NewUserRepository(db).FindUserByEmail(email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindUserByEmail_QueryError(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	email := "john.doe@example.com"

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."email" = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnError(errors.New("query error"))

	user, err := NewUserRepository(db).FindUserByEmail(email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "query error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_CreateUser_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	user := utils.GenerateRandomUser()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(user.ID, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx := db.Begin()
	err := NewUserRepository(db).CreateUser(tx, user)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_CreateUser_Failure(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	user := utils.GenerateRandomUser()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(user.ID, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	tx := db.Begin()
	err := NewUserRepository(db).CreateUser(tx, user)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.Error(t, err)
	assert.Equal(t, "insert error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
