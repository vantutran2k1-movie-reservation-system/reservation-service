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
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserRepository(db)

	user := utils.GenerateUser()
	profile := utils.GenerateUserProfile()
	profile.UserID = user.ID
	user.Profile = profile
	filter := filters.UserFilter{
		Filter: &filters.SingleFilter{},
		ID:     &filters.Condition{Operator: filters.OpEqual, Value: user.ID},
		Email:  &filters.Condition{Operator: filters.OpEqual, Value: user.Email},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND email = $2 ORDER BY "users"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Email.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(user))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_profiles" WHERE "user_profiles"."user_id" = $1`)).
			WithArgs(user.ID).
			WillReturnRows(utils.GenerateSqlMockRow(profile))

		result, err := repo.GetUser(filter, true)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND email = $2 ORDER BY "users"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Email.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetUser(filter, true)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting user", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND email = $2 ORDER BY "users"."id" LIMIT $3`)).
			WithArgs(filter.ID.Value, filter.Email.Value, 1).
			WillReturnError(errors.New("error getting user"))

		result, err := repo.GetUser(filter, true)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "error getting user", err.Error())
	})
}

func TestUserRepository_UserExists(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserRepository(db)

	user := utils.GenerateUser()
	filter := filters.UserFilter{
		Filter: &filters.SingleFilter{},
		Email:  &filters.Condition{Operator: filters.OpEqual, Value: user.Email},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(filter.Email.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(user))

		result, err := repo.UserExists(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(filter.Email.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.UserExists(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, false, result)
	})

	t.Run("error getting user", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id" FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(filter.Email.Value, 1).
			WillReturnError(errors.New("error getting user"))

		result, err := repo.UserExists(filter)

		assert.NotNil(t, result)
		assert.Error(t, err)
		assert.Equal(t, false, result)
		assert.Equal(t, "error getting user", err.Error())
	})
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	user := utils.GenerateUser()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash","is_active","is_verified","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`)).
			WithArgs(user.ID, user.Email, user.PasswordHash, user.IsActive, user.IsVerified, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := NewUserRepository(db).CreateUser(tx, user)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash","is_active","is_verified","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`)).
			WithArgs(user.ID, user.Email, user.PasswordHash, user.IsActive, user.IsVerified, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := NewUserRepository(db).CreateUser(tx, user)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestUserRepository_CreateOrUpdateUser(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock.ExpectClose()
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserRepository(db)

	user := utils.GenerateUser()
	newUser := utils.GenerateUser()
	newUser.Email = user.Email

	t.Run("success creating new user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(newUser.Email, 1).
			WillReturnRows(sqlmock.NewRows(nil))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash","is_active","is_verified","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`)).
			WithArgs(newUser.ID, newUser.Email, newUser.PasswordHash, newUser.IsActive, newUser.IsVerified, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		result, err := repo.CreateOrUpdateUser(tx, newUser)
		tx.Commit()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		//assert.Equal(t, newUser.ID, result.ID)
		//assert.Equal(t, newUser.Email, result.Email)
		//assert.Equal(t, newUser.PasswordHash, result.PasswordHash)
	})

	t.Run("success updating existing user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(newUser.Email, 1).
			WillReturnRows(utils.GenerateSqlMockRow(user))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "email"=$1,"password_hash"=$2,"is_active"=$3,"is_verified"=$4,"created_at"=$5,"updated_at"=$6 WHERE "id" = $7`)).
			WithArgs(newUser.Email, newUser.PasswordHash, newUser.IsActive, newUser.IsVerified, sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		result, err := repo.CreateOrUpdateUser(tx, newUser)
		tx.Commit()

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, newUser.Email, result.Email)
		assert.Equal(t, newUser.PasswordHash, result.PasswordHash)
	})

	t.Run("error getting user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(newUser.Email, 1).
			WillReturnError(errors.New("error getting user"))
		mock.ExpectRollback()

		tx := db.Begin()
		result, err := repo.CreateOrUpdateUser(tx, newUser)
		tx.Rollback()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting user", err.Error())
	})

	t.Run("error creating new user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(newUser.Email, 1).
			WillReturnRows(sqlmock.NewRows(nil))
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","email","password_hash","is_active","is_verified","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7)`)).
			WithArgs(newUser.ID, newUser.Email, newUser.PasswordHash, newUser.IsActive, newUser.IsVerified, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("error creating new user"))
		mock.ExpectRollback()

		tx := db.Begin()
		result, err := repo.CreateOrUpdateUser(tx, newUser)
		tx.Rollback()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error creating new user", err.Error())
	})

	t.Run("error updating existing user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(newUser.Email, 1).
			WillReturnRows(utils.GenerateSqlMockRow(user))
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "email"=$1,"password_hash"=$2,"is_active"=$3,"is_verified"=$4,"created_at"=$5,"updated_at"=$6 WHERE "id" = $7`)).
			WithArgs(newUser.Email, newUser.PasswordHash, newUser.IsActive, newUser.IsVerified, sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID).
			WillReturnError(errors.New("error updating existing user"))
		mock.ExpectRollback()

		tx := db.Begin()
		result, err := repo.CreateOrUpdateUser(tx, newUser)
		tx.Rollback()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error updating existing user", err.Error())
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
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

func TestUserRepository_VerifyUser(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.Nil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserRepository(db)

	user := utils.GenerateUser()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "is_verified"=$1,"updated_at"=$2 WHERE "id" = $3`)).
			WithArgs(true, sqlmock.AnyArg(), user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.VerifyUser(tx, user)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("error updating user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "is_verified"=$1,"updated_at"=$2 WHERE "id" = $3`)).
			WithArgs(true, sqlmock.AnyArg(), user.ID).
			WillReturnError(errors.New("error updating user"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.VerifyUser(tx, user)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error updating user")
	})
}
