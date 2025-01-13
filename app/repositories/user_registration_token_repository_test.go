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

func TestUserRegistrationTokenRepository_GetToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserRegistrationTokenRepository(db)

	token := utils.GenerateUserRegistrationToken()
	filter := filters.UserRegistrationTokenFilter{
		Filter:     &filters.SingleFilter{},
		TokenValue: &filters.Condition{Operator: filters.OpEqual, Value: token.TokenValue},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_registration_tokens" WHERE token_value = $1 ORDER BY "user_registration_tokens"."id" LIMIT $2`)).
			WithArgs(token.TokenValue, 1).
			WillReturnRows(utils.GenerateSqlMockRow(token))

		result, err := repo.GetToken(filter)

		assert.NotNil(t, result)
		assert.NoError(t, err)
		assert.Equal(t, token, result)
	})

	t.Run("token not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_registration_tokens" WHERE token_value = $1 ORDER BY "user_registration_tokens"."id" LIMIT $2`)).
			WithArgs(token.TokenValue, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetToken(filter)

		assert.Nil(t, result)
		assert.NoError(t, err)
	})

	t.Run("error getting token", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_registration_tokens" WHERE token_value = $1 ORDER BY "user_registration_tokens"."id" LIMIT $2`)).
			WithArgs(token.TokenValue, 1).
			WillReturnError(errors.New("error getting token"))

		result, err := repo.GetToken(filter)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "error getting token")
	})
}

func TestUserRegistrationTokenRepository_CreateToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserRegistrationTokenRepository(db)

	token := utils.GenerateUserRegistrationToken()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "user_registration_tokens" ("id","user_id","token_value","is_used","created_at","expires_at") VALUES ($1,$2,$3,$4,$5,$6)`)).
			WithArgs(token.ID, token.UserID, token.TokenValue, token.IsUsed, token.CreatedAt, token.ExpiresAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateToken(tx, token)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "user_registration_tokens" ("id","user_id","token_value","is_used","created_at","expires_at") VALUES ($1,$2,$3,$4,$5,$6)`)).
			WithArgs(token.ID, token.UserID, token.TokenValue, token.IsUsed, token.CreatedAt, token.ExpiresAt).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateToken(tx, token)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "db error")
	})
}

func TestUserRegistrationTokenRepository_UseToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewUserRegistrationTokenRepository(db)

	token := utils.GenerateUserRegistrationToken()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_registration_tokens" SET "is_used"=$1 WHERE "id" = $2`)).
			WithArgs(true, token.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.UseToken(tx, token)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_registration_tokens" SET "is_used"=$1 WHERE "id" = $2`)).
			WithArgs(true, token.ID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UseToken(tx, token)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.EqualError(t, err, "db error")
	})
}
