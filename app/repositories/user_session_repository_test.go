package repositories

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestUserSessionRepository_GetUserSession_Success(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionID := "test-session-id"
	session := utils.GenerateRandomUserSession()

	sessionJSON, err := json.Marshal(session)
	require.NoError(t, err)

	mock.ExpectGet(sessionID).SetVal(string(sessionJSON))

	result, err := NewUserSessionRepository(redisClient).GetUserSession(sessionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, session.UserID, result.UserID)
	assert.Equal(t, session.Email, result.Email)
}

func TestUserSessionRepository_GetUserSession_NotFound(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionID := "non-existent-session-id"

	mock.ExpectGet(sessionID).RedisNil()

	result, err := NewUserSessionRepository(redisClient).GetUserSession(sessionID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserSessionRepository_GetUserSession_JsonUnmarshalError(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionID := "test-session-id"
	mock.ExpectGet(sessionID).SetVal("invalid json")

	result, err := NewUserSessionRepository(redisClient).GetUserSession(sessionID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserSessionRepository_CreateUserSession_Success(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionID := "test-session-id"
	expiration := 24 * time.Hour
	session := utils.GenerateRandomUserSession()

	sessionJSON, err := json.Marshal(session)
	require.NoError(t, err)

	mock.ExpectSet(sessionID, sessionJSON, expiration).SetVal("OK")

	err = NewUserSessionRepository(redisClient).CreateUserSession(sessionID, expiration, session)

	assert.NoError(t, err)
}

func TestUserSessionRepository_CreateUserSession_MarshalError(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionID := "test-session-id"
	expiration := 24 * time.Hour
	session := utils.GenerateRandomUserSession()

	session.Email = string([]byte{255})

	err := NewUserSessionRepository(redisClient).CreateUserSession(sessionID, expiration, session)

	assert.Error(t, err)
}

func TestUserSessionRepository_CreateUserSession_SetError(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionID := "test-session-id"
	expiration := 24 * time.Hour
	session := utils.GenerateRandomUserSession()

	sessionJSON, err := json.Marshal(session)
	require.NoError(t, err)

	mock.ExpectSet(sessionID, sessionJSON, expiration).SetErr(redis.Nil)

	err = NewUserSessionRepository(redisClient).CreateUserSession(sessionID, expiration, session)

	assert.Error(t, err)
}

func TestUserSessionRepository_DeleteUserSession_Success(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionID := "test-session-id"

	mock.ExpectDel(sessionID).SetVal(1)

	err := NewUserSessionRepository(redisClient).DeleteUserSession(sessionID)

	assert.NoError(t, err)
}

func TestUserSessionRepository_DeleteUserSession_Error(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionID := "test-session-id"

	mock.ExpectDel(sessionID).SetErr(redis.Nil)

	err := NewUserSessionRepository(redisClient).DeleteUserSession(sessionID)

	assert.Error(t, err)
}

func TestUserSessionRepository_DeleteUserSessions_Success(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionKey := "session-key"
	session := utils.GenerateRandomUserSession()
	sessionData, _ := json.Marshal(session)

	mock.ExpectScan(0, "*", 100).SetVal([]string{sessionKey}, 0)
	mock.ExpectGet(sessionKey).SetVal(string(sessionData))
	mock.ExpectDel(sessionKey).SetVal(1)

	err := NewUserSessionRepository(redisClient).DeleteUserSessions(session.UserID)

	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func TestUserSessionRepository_DeleteUserSessions_ErrorScanningKeys(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	mock.ExpectScan(0, "*", 100).SetErr(errors.New("scan error"))

	err := NewUserSessionRepository(redisClient).DeleteUserSessions(uuid.New())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error scanning keys")
}

func TestUserSessionRepository_DeleteUserSessions_ErrorGettingKey(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	userID := uuid.New()
	sessionKey := "session-key"
	mock.ExpectScan(0, "*", 100).SetVal([]string{sessionKey}, 0)
	mock.ExpectGet(sessionKey).SetErr(errors.New("get error"))

	err := NewUserSessionRepository(redisClient).DeleteUserSessions(userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting key")
}

func TestUserSessionRepository_DeleteUserSessions_ErrorDeletingKey(t *testing.T) {
	redisClient, mock := mock_db.SetupTestRedis()
	defer func() {
		require.NoError(t, mock_db.TearDownTestRedis(mock))
	}()

	sessionKey := "session-key"
	session := utils.GenerateRandomUserSession()
	sessionData, _ := json.Marshal(session)

	mock.ExpectScan(0, "*", 100).SetVal([]string{sessionKey}, 0)
	mock.ExpectGet(sessionKey).SetVal(string(sessionData))
	mock.ExpectDel(sessionKey).SetErr(errors.New("delete error"))

	err := NewUserSessionRepository(redisClient).DeleteUserSessions(session.UserID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting key")
}
