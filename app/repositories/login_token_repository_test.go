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
	"time"
)

func TestLoginTokenRepository_GetLoginToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewLoginTokenRepository(db)

	token := utils.GenerateLoginToken()
	filter := filters.LoginTokenFilter{
		Filter:     &filters.SingleFilter{Logic: filters.And},
		UserID:     &filters.Condition{Operator: filters.OpEqual, Value: token.UserID},
		TokenValue: &filters.Condition{Operator: filters.OpEqual, Value: token.TokenValue},
		ExpiresAt:  &filters.Condition{Operator: filters.OpGreater, Value: time.Now().UTC()},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "login_tokens" WHERE user_id = $1 AND token_value = $2 AND expires_at > $3 ORDER BY "login_tokens"."id" LIMIT $4`)).
			WithArgs(filter.UserID.Value, filter.TokenValue.Value, filter.ExpiresAt.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(token))

		result, err := repo.GetLoginToken(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, token, result)
	})

	t.Run("token not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "login_tokens" WHERE user_id = $1 AND token_value = $2 AND expires_at > $3 ORDER BY "login_tokens"."id" LIMIT $4`)).
			WithArgs(filter.UserID.Value, filter.TokenValue.Value, filter.ExpiresAt.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetLoginToken(filter)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting token", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "login_tokens" WHERE user_id = $1 AND token_value = $2 AND expires_at > $3 ORDER BY "login_tokens"."id" LIMIT $4`)).
			WithArgs(filter.UserID.Value, filter.TokenValue.Value, filter.ExpiresAt.Value, 1).
			WillReturnError(errors.New("error getting token"))

		result, err := NewLoginTokenRepository(db).GetLoginToken(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting token", err.Error())
	})
}

func TestLoginTokenRepository_CreateLoginToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewLoginTokenRepository(db)

	token := utils.GenerateLoginToken()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "login_tokens" ("id","user_id","token_value","created_at","expires_at") VALUES ($1,$2,$3,$4,$5)`)).
			WithArgs(token.ID, token.UserID, token.TokenValue, token.CreatedAt, token.ExpiresAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.CreateLoginToken(tx, token)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "login_tokens" ("id","user_id","token_value","created_at","expires_at") VALUES ($1,$2,$3,$4,$5)`)).
			WithArgs(token.ID, token.UserID, token.TokenValue, token.CreatedAt, token.ExpiresAt).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateLoginToken(tx, token)
		tx.Rollback()

		assert.NotNil(t, err)
	})
}

func TestLoginTokenRepository_RevokeLoginToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewLoginTokenRepository(db)

	token := utils.GenerateLoginToken()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "login_tokens" SET "expires_at"=$1 WHERE "id" = $2`)).
			WithArgs(sqlmock.AnyArg(), token.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.RevokeLoginToken(tx, token)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "login_tokens" SET "expires_at"=$1 WHERE "id" = $2`)).
			WithArgs(sqlmock.AnyArg(), token.ID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.RevokeLoginToken(tx, token)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestLoginTokenRepository_RevokeUserLoginTokens(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewLoginTokenRepository(db)

	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "login_tokens" SET "expires_at"=$1 WHERE user_id = $2 AND expires_at > $3`)).
			WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.RevokeUserLoginTokens(tx, userID)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "login_tokens" SET "expires_at"=$1 WHERE user_id = $2 AND expires_at > $3`)).
			WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.RevokeUserLoginTokens(tx, userID)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
