package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
			allowed, ttl := service.Allow(clientIp)
			fmt.Print(allowed)
			fmt.Print(ttl)
		}
	})
}
