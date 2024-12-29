package repositories

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

func TestUserSessionRepository_GetUserSession(t *testing.T) {
	client, mock := mock_db.SetupTestRedis()
	defer func() {
		assert.Nil(t, mock_db.TearDownTestRedis(mock))
	}()

	repo := NewUserSessionRepository(client)

	session := utils.GenerateUserSession()
	sessionID := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		sessionJSON, _ := json.Marshal(session)
		mock.ExpectGet(sessionID).SetVal(string(sessionJSON))

		result, err := repo.GetUserSession(sessionID)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, session.UserID, result.UserID)
		assert.Equal(t, session.Email, result.Email)
	})

	t.Run("session not found", func(t *testing.T) {
		mock.ExpectGet(sessionID).RedisNil()

		result, err := repo.GetUserSession(sessionID)

		assert.Nil(t, result)
		assert.Nil(t, err)
	})

	t.Run("error unmarshalling data", func(t *testing.T) {
		mock.ExpectGet(sessionID).SetVal("invalid json")

		result, err := repo.GetUserSession(sessionID)

		assert.Nil(t, result)
		assert.NotNil(t, err)
	})
}

func TestUserSessionRepository_CreateUserSession(t *testing.T) {
	client, mock := mock_db.SetupTestRedis()
	defer func() {
		assert.Nil(t, mock_db.TearDownTestRedis(mock))
	}()

	repo := NewUserSessionRepository(client)

	session := utils.GenerateUserSession()
	sessionID := uuid.NewString()
	expiration := 24 * time.Hour

	t.Run("success", func(t *testing.T) {
		sessionJSON, _ := json.Marshal(session)
		mock.ExpectSet(sessionID, sessionJSON, expiration).SetVal("OK")

		err := repo.CreateUserSession(sessionID, expiration, session)

		assert.Nil(t, err)
	})

	t.Run("error marshalling data", func(t *testing.T) {
		session.Email = string([]byte{255})

		err := repo.CreateUserSession(sessionID, expiration, session)

		assert.NotNil(t, err)
	})

	t.Run("error creating session", func(t *testing.T) {
		sessionJSON, _ := json.Marshal(session)
		mock.ExpectSet(sessionID, sessionJSON, expiration).SetErr(redis.Nil)

		err := repo.CreateUserSession(sessionID, expiration, session)

		assert.NotNil(t, err)
	})
}

func TestUserSessionRepository_DeleteUserSession(t *testing.T) {
	client, mock := mock_db.SetupTestRedis()
	defer func() {
		assert.Nil(t, mock_db.TearDownTestRedis(mock))
	}()

	repo := NewUserSessionRepository(client)

	sessionID := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		mock.ExpectDel(sessionID).SetVal(1)

		err := repo.DeleteUserSession(sessionID)

		assert.Nil(t, err)
	})

	t.Run("error deleting session", func(t *testing.T) {
		mock.ExpectDel(sessionID).SetErr(redis.Nil)

		err := repo.DeleteUserSession(sessionID)

		assert.NotNil(t, err)
	})
}

func TestUserSessionRepository_DeleteUserSessions(t *testing.T) {
	client, mock := mock_db.SetupTestRedis()
	defer func() {
		assert.Nil(t, mock_db.TearDownTestRedis(mock))
	}()

	repo := NewUserSessionRepository(client)

	session := utils.GenerateUserSession()
	sessionID := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		sessionData, _ := json.Marshal(session)

		mock.ExpectScan(0, "*", 100).SetVal([]string{sessionID}, 0)
		mock.ExpectGet(sessionID).SetVal(string(sessionData))
		mock.ExpectDel(sessionID).SetVal(1)

		err := repo.DeleteUserSessions(session.UserID)

		assert.Nil(t, err)
	})

	t.Run("error scanning keys", func(t *testing.T) {
		mock.ExpectScan(0, "*", 100).SetErr(errors.New("error scanning keys"))

		err := repo.DeleteUserSessions(uuid.New())

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error scanning keys")
	})

	t.Run("error getting keys", func(t *testing.T) {
		mock.ExpectScan(0, "*", 100).SetVal([]string{sessionID}, 0)
		mock.ExpectGet(sessionID).SetErr(errors.New("error getting keys"))

		err := repo.DeleteUserSessions(session.UserID)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error getting keys")
	})

	t.Run("error deleting sessions", func(t *testing.T) {
		sessionData, _ := json.Marshal(session)

		mock.ExpectScan(0, "*", 100).SetVal([]string{sessionID}, 0)
		mock.ExpectGet(sessionID).SetVal(string(sessionData))
		mock.ExpectDel(sessionID).SetErr(errors.New("error deleting sessions"))

		err := repo.DeleteUserSessions(session.UserID)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error deleting sessions")
	})
}
