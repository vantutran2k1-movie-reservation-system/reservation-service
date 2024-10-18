package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"regexp"
	"testing"
)

func TestPasswordResetTokenRepository_GetActivePasswordResetToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewPasswordResetTokenRepository(db)

	token := utils.GeneratePasswordResetToken()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "token_value", "is_used", "created_at", "expires_at"}).
			AddRow(token.ID, token.UserID, token.TokenValue, token.IsUsed, token.CreatedAt, token.ExpiresAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE token_value = $1 AND is_used = $2 AND expires_at > $3 ORDER BY "password_reset_tokens"."id" LIMIT $4`)).
			WithArgs(token.TokenValue, false, sqlmock.AnyArg(), 1).
			WillReturnRows(rows)

		result, err := repo.GetActivePasswordResetToken(token.TokenValue)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, token, result)
	})

	t.Run("token not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE token_value = $1 AND is_used = $2 AND expires_at > $3 ORDER BY "password_reset_tokens"."id" LIMIT $4`)).
			WithArgs(token.TokenValue, false, sqlmock.AnyArg(), 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetActivePasswordResetToken(token.TokenValue)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE token_value = $1 AND is_used = $2 AND expires_at > $3 ORDER BY "password_reset_tokens"."id" LIMIT $4`)).
			WithArgs(token.TokenValue, false, sqlmock.AnyArg(), 1).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetActivePasswordResetToken(token.TokenValue)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestPasswordResetTokenRepository_GetUserActivePasswordResetTokens(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		assert.NotNil(t, mock_db.TearDownTestDB(db, mock))
	}()

	repo := NewPasswordResetTokenRepository(db)

	user := utils.GenerateUser()

	t.Run("success", func(t *testing.T) {
		tokens := utils.GeneratePasswordResetTokens(3)

		rows := sqlmock.NewRows([]string{"id", "user_id", "token_value", "is_used", "created_at", "expires_at"})
		for _, token := range tokens {
			rows.AddRow(token.ID, token.UserID, token.TokenValue, token.IsUsed, token.CreatedAt, token.ExpiresAt)
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE user_id = $1 AND is_used = $2 AND expires_at > $3`)).
			WithArgs(user.ID, false, sqlmock.AnyArg()).
			WillReturnRows(rows)

		result, err := repo.GetUserActivePasswordResetTokens(user.ID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, tokens, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "password_reset_tokens" WHERE user_id = $1 AND is_used = $2 AND expires_at > $3`)).
			WithArgs(user.ID, false, sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetUserActivePasswordResetTokens(user.ID)

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
