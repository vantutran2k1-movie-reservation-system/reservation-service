package services

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/mocks/mock_db"
	"testing"
	"time"
)

func TestRateLimiterService_Allow(t *testing.T) {
	client, mock := mock_db.SetupTestRedis()
	defer func() {
		assert.Nil(t, mock_db.TearDownTestRedis(mock))
	}()

	maxRequests := 5
	service := NewRateLimiterService(client, maxRequests, time.Minute)

	clientIp := "127.0.0.1"

	t.Run("should allow within limit", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			mock.ExpectIncr(fmt.Sprintf("%s:%s", constants.ClientRateLimit, clientIp)).SetVal(int64(i))

			allowed, ttl := service.Allow(clientIp)

			assert.True(t, allowed)
			assert.Equal(t, time.Duration(0), ttl)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectIncr(fmt.Sprintf("%s:%s", constants.ClientRateLimit, clientIp)).SetErr(errors.New("db error"))

		allowed, ttl := service.Allow(clientIp)

		assert.False(t, allowed)
		assert.Equal(t, time.Duration(0), ttl)
	})

	t.Run("should not allow after hitting limits", func(t *testing.T) {
		mock.ExpectIncr(fmt.Sprintf("%s:%s", constants.ClientRateLimit, clientIp)).SetVal(int64(maxRequests + 1))
		allowed, ttl := service.Allow(clientIp)

		assert.False(t, allowed)
		assert.Equal(t, time.Duration(0), ttl)
	})
}
