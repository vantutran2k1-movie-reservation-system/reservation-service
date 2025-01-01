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
