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
	"time"
)

func TestPasswordResetTokenRepository_GetToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewPasswordResetTokenRepository(db)

	token := utils.GeneratePasswordResetToken()
	filter := filters.PasswordResetTokenFilter{
		Filter:     &filters.SingleFilter{},
		TokenValue: &filters.Condition{Operator: filters.OpEqual, Value: token.TokenValue},
		IsUsed:     &filters.Condition{Operator: filters.OpEqual, Value: false},
		ExpiresAt:  &filters.Condition{Operator: filters.OpGreater, Value: time.Now().UTC()},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE token_value = $1 AND is_used = $2 AND expires_at > $3 ORDER BY "password_reset_tokens"."id" LIMIT $4`)).
			WithArgs(filter.TokenValue.Value, filter.IsUsed.Value, filter.ExpiresAt.Value, 1).
			WillReturnRows(utils.GenerateSqlMockRow(token))

		result, err := repo.GetToken(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, token, result)
	})

	t.Run("token not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE token_value = $1 AND is_used = $2 AND expires_at > $3 ORDER BY "password_reset_tokens"."id" LIMIT $4`)).
			WithArgs(filter.TokenValue.Value, filter.IsUsed.Value, filter.ExpiresAt.Value, 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetToken(filter)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error getting token", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE token_value = $1 AND is_used = $2 AND expires_at > $3 ORDER BY "password_reset_tokens"."id" LIMIT $4`)).
			WithArgs(filter.TokenValue.Value, filter.IsUsed.Value, filter.ExpiresAt.Value, 1).
			WillReturnError(errors.New("error getting token"))

		result, err := repo.GetToken(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting token", err.Error())
	})
}

func TestPasswordResetTokenRepository_GetTokens(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewPasswordResetTokenRepository(db)

	user := utils.GenerateUser()
	filter := filters.PasswordResetTokenFilter{
		Filter:    &filters.MultiFilter{},
		UserID:    &filters.Condition{Operator: filters.OpEqual, Value: user.ID},
		IsUsed:    &filters.Condition{Operator: filters.OpEqual, Value: false},
		ExpiresAt: &filters.Condition{Operator: filters.OpGreater, Value: time.Now().UTC()},
	}

	t.Run("success", func(t *testing.T) {
		tokens := utils.GeneratePasswordResetTokens(3)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE user_id = $1 AND is_used = $2 AND expires_at > $3`)).
			WithArgs(filter.UserID.Value, filter.IsUsed.Value, filter.ExpiresAt.Value).
			WillReturnRows(utils.GenerateSqlMockRows(tokens))

		result, err := repo.GetTokens(filter)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, tokens, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE user_id = $1 AND is_used = $2 AND expires_at > $3`)).
			WithArgs(filter.UserID.Value, filter.IsUsed.Value, filter.ExpiresAt.Value).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetTokens(filter)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPasswordResetTokenRepository_CreateToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewPasswordResetTokenRepository(db)

	token := utils.GeneratePasswordResetToken()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "password_reset_tokens" ("id","user_id","token_value","is_used","created_at","expires_at") VALUES ($1,$2,$3,$4,$5,$6)`)).
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
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "password_reset_tokens" ("id","user_id","token_value","is_used","created_at","expires_at") VALUES ($1,$2,$3,$4,$5,$6)`)).
			WithArgs(token.ID, token.UserID, token.TokenValue, token.IsUsed, token.CreatedAt, token.ExpiresAt).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.CreateToken(tx, token)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestPasswordResetTokenRepository_UseToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewPasswordResetTokenRepository(db)

	token := utils.GeneratePasswordResetToken()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "password_reset_tokens" SET "is_used"=$1 WHERE "id" = $2`)).
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
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "password_reset_tokens" SET "is_used"=$1 WHERE "id" = $2`)).
			WithArgs(true, token.ID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.UseToken(tx, token)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestPasswordResetTokenRepository_RevokeTokens(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewPasswordResetTokenRepository(db)

	tokens := utils.GeneratePasswordResetTokens(3)

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "password_reset_tokens" SET "expires_at"=$1 WHERE id IN ($2,$3,$4)`)).
			WithArgs(sqlmock.AnyArg(), tokens[0].ID, tokens[1].ID, tokens[2].ID).
			WillReturnResult(sqlmock.NewResult(3, 3))
		mock.ExpectCommit()

		tx := db.Begin()
		err := repo.RevokeTokens(tx, tokens)
		tx.Commit()

		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "password_reset_tokens" SET "expires_at"=$1 WHERE id IN ($2,$3,$4)`)).
			WithArgs(sqlmock.AnyArg(), tokens[0].ID, tokens[1].ID, tokens[2].ID).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		tx := db.Begin()
		err := repo.RevokeTokens(tx, tokens)
		tx.Rollback()

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
