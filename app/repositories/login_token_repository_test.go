package repositories

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/test"
	"gorm.io/gorm"
)

func TestLoginTokenRepository_GetActiveLoginToken_Success(t *testing.T) {
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
	}()

	expectedToken := test.GenerateRandomLoginToken()

	mock.ExpectQuery(`SELECT \* FROM "login_tokens" WHERE token_value = \$1 AND expires_at > \$2 ORDER BY "login_tokens"."id" LIMIT \$3`).
		WithArgs(expectedToken.TokenValue, sqlmock.AnyArg(), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token_value", "created_at", "expires_at"}).
			AddRow(expectedToken.ID, expectedToken.UserID, expectedToken.TokenValue, expectedToken.CreatedAt, expectedToken.ExpiresAt))

	result, err := NewLoginTokenRepository(db).GetActiveLoginToken(expectedToken.TokenValue)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedToken.ID, result.ID)
	assert.Equal(t, expectedToken.UserID, result.UserID)
	assert.Equal(t, expectedToken.TokenValue, result.TokenValue)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_GetActiveLoginToken_Failure_NotFound(t *testing.T) {
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
	}()

	tokenValue := "non_existent_token"

	mock.ExpectQuery(`SELECT \* FROM "login_tokens" WHERE token_value = \$1 AND expires_at > \$2 ORDER BY "login_tokens"."id" LIMIT \$3`).
		WithArgs(tokenValue, sqlmock.AnyArg(), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := NewLoginTokenRepository(db).GetActiveLoginToken(tokenValue)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_GetActiveLoginToken_Failure_ExpiredToken(t *testing.T) {
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
	}()

	tokenValue := "expired_token"

	mock.ExpectQuery(`SELECT \* FROM "login_tokens" WHERE token_value = \$1 AND expires_at > \$2 ORDER BY "login_tokens"."id" LIMIT \$3`).
		WithArgs(tokenValue, sqlmock.AnyArg(), 1).
		WillReturnRows(sqlmock.NewRows(nil))

	result, err := NewLoginTokenRepository(db).GetActiveLoginToken(tokenValue)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginTokenRepository_CreateLoginToken_Success(t *testing.T) {
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
	}()

	loginToken := test.GenerateRandomLoginToken()

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
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
	}()

	loginToken := test.GenerateRandomLoginToken()

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
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
	}()

	loginToken := test.GenerateRandomLoginToken()

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
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
	}()

	loginToken := test.GenerateRandomLoginToken()

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
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
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
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
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
	db, mock := test.SetupTestDB(t)
	defer func() {
		test.TearDownTestDB(db, mock)
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
