package repositories

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"gorm.io/gorm"
)

func TestLoginTokenRepository_GetActiveLoginToken(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	repo := NewLoginTokenRepository(db)

	token := utils.GenerateRandomLoginToken()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "token_value", "created_at", "expires_at"}).
			AddRow(token.ID, token.UserID, token.TokenValue, token.CreatedAt, token.ExpiresAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "login_tokens" WHERE token_value = $1 AND expires_at > $2 ORDER BY "login_tokens"."id" LIMIT $3`)).
			WithArgs(token.TokenValue, sqlmock.AnyArg(), 1).
			WillReturnRows(rows)

		result, err := repo.GetActiveLoginToken(token.TokenValue)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, token, result)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("token not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "login_tokens" WHERE token_value = $1 AND expires_at > $2 ORDER BY "login_tokens"."id" LIMIT $3`)).
			WithArgs(token.TokenValue, sqlmock.AnyArg(), 1).
			WillReturnRows(sqlmock.NewRows(nil))

		result, err := repo.GetActiveLoginToken(token.TokenValue)

		assert.Nil(t, result)
		assert.Nil(t, err)

		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "login_tokens" WHERE token_value = $1 AND expires_at > $2 ORDER BY "login_tokens"."id" LIMIT $3`)).
			WithArgs(token.TokenValue, sqlmock.AnyArg(), 1).
			WillReturnError(errors.New("db error"))

		result, err := NewLoginTokenRepository(db).GetActiveLoginToken(token.TokenValue)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestLoginTokenRepository_CreateLoginToken_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	loginToken := utils.GenerateRandomLoginToken()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "login_tokens"`).
		WithArgs(loginToken.ID, loginToken.UserID, loginToken.TokenValue, loginToken.CreatedAt, loginToken.ExpiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx := db.Begin()
	err := NewLoginTokenRepository(db).CreateLoginToken(tx, loginToken)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_CreateLoginToken_Failure(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	loginToken := utils.GenerateRandomLoginToken()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "login_tokens"`).
		WithArgs(loginToken.ID, loginToken.UserID, loginToken.TokenValue, loginToken.CreatedAt, loginToken.ExpiresAt).
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	tx := db.Begin()
	err := NewLoginTokenRepository(db).CreateLoginToken(tx, loginToken)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidData, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_RevokeLoginToken_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	loginToken := utils.GenerateRandomLoginToken()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "login_tokens" SET "user_id"=\$1,"token_value"=\$2,"created_at"=\$3,"expires_at"=\$4 WHERE "id" = \$5`).
		WithArgs(loginToken.UserID, loginToken.TokenValue, sqlmock.AnyArg(), sqlmock.AnyArg(), loginToken.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx := db.Begin()
	err := NewLoginTokenRepository(db).RevokeLoginToken(tx, loginToken)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().UTC(), loginToken.ExpiresAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_RevokeLoginToken_Failure(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	loginToken := utils.GenerateRandomLoginToken()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "login_tokens" SET "user_id"=\$1,"token_value"=\$2,"created_at"=\$3,"expires_at"=\$4 WHERE "id" = \$5`).
		WithArgs(loginToken.UserID, loginToken.TokenValue, sqlmock.AnyArg(), sqlmock.AnyArg(), loginToken.ID).
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	tx := db.Begin()
	err := NewLoginTokenRepository(db).RevokeLoginToken(tx, loginToken)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidData, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_RevokeUserLoginTokens_Success(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	userID := uuid.New()
	now := time.Now().UTC()

	tokens := []models.LoginToken{
		{
			ID:         uuid.New(),
			UserID:     userID,
			TokenValue: "token1",
			CreatedAt:  now,
			ExpiresAt:  now.Add(24 * time.Hour),
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			TokenValue: "token2",
			CreatedAt:  now,
			ExpiresAt:  now.Add(24 * time.Hour),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "token_value", "created_at", "expires_at"}).
		AddRow(tokens[0].ID, tokens[0].UserID, tokens[0].TokenValue, tokens[0].CreatedAt, tokens[0].ExpiresAt).
		AddRow(tokens[1].ID, tokens[1].UserID, tokens[1].TokenValue, tokens[1].CreatedAt, tokens[1].ExpiresAt)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT \* FROM "login_tokens" WHERE user_id = \$1 AND expires_at > \$2`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(rows)

	mock.ExpectExec(`INSERT INTO "login_tokens" \("id","user_id","token_value","created_at","expires_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\),\(\$6,\$7,\$8,\$9,\$10\) ON CONFLICT \("id"\) DO UPDATE SET "user_id"="excluded"\."user_id","token_value"="excluded"\."token_value","expires_at"="excluded"\."expires_at"`).
		WithArgs(tokens[0].ID, tokens[0].UserID, tokens[0].TokenValue, sqlmock.AnyArg(), sqlmock.AnyArg(), tokens[1].ID, tokens[1].UserID, tokens[1].TokenValue, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(2, 2))
	mock.ExpectCommit()

	tx := db.Begin()
	err := NewLoginTokenRepository(db).RevokeUserLoginTokens(tx, userID)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_RevokeUserLoginTokens_Failure_QueryError(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	userID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT \* FROM "login_tokens"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnError(errors.New("failed to retrieve tokens"))
	mock.ExpectRollback()

	tx := db.Begin()
	err := NewLoginTokenRepository(db).RevokeUserLoginTokens(tx, userID)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to retrieve tokens")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_RevokeUserLoginTokens_Failure_UpdateError(t *testing.T) {
	db, mock := mock_db.SetupTestDB(t)
	defer func() {
		mock_db.TearDownTestDB(db, mock)
	}()

	userID := uuid.New()
	now := time.Now().UTC()

	rows := sqlmock.NewRows([]string{"id", "user_id", "token_value", "created_at", "expires_at"}).
		AddRow(uuid.New(), userID, "token1", now, now.Add(24*time.Hour)).
		AddRow(uuid.New(), userID, "token2", now, now.Add(24*time.Hour))

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT \* FROM "login_tokens"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(rows)
	mock.ExpectExec(`INSERT INTO "login_tokens" \("id","user_id","token_value","created_at","expires_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\),\(\$6,\$7,\$8,\$9,\$10\) ON CONFLICT \("id"\) DO UPDATE SET "user_id"="excluded"\."user_id","token_value"="excluded"\."token_value","expires_at"="excluded"\."expires_at"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("failed to update tokens"))
	mock.ExpectRollback()

	tx := db.Begin()
	err := NewLoginTokenRepository(db).RevokeUserLoginTokens(tx, userID)
	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to update tokens")
	assert.NoError(t, mock.ExpectationsWereMet())
}
