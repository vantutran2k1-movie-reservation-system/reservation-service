package services

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

type RateLimiterService interface {
	Allow(clientIp string) (bool, time.Duration)
}

func NewRateLimiterService(
	redisClient *redis.Client,
	maxRequests int,
	windowTime time.Duration,
) RateLimiterService {
	return &rateLimiterService{
		redisClient: redisClient,
		maxRequests: maxRequests,
		windowTime:  windowTime,
	}
}

type rateLimiterService struct {
	redisClient *redis.Client
	maxRequests int
	windowTime  time.Duration
}

func (s *rateLimiterService) Allow(clientIp string) (bool, time.Duration) {
	redisKey := fmt.Sprintf("rate_limit:%s", clientIp)
	
	count, err := s.redisClient.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, 0
	}

	if count == 1 {
		s.redisClient.Expire(ctx, redisKey, s.windowTime)
	}

	if count > int64(s.maxRequests) {
		ttl, _ := s.redisClient.TTL(ctx, redisKey).Result()
		return false, ttl
	}

	return true, 0
}
