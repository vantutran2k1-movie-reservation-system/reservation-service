package mocks

import (
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func SetupTestRedis() (*redis.Client, redismock.ClientMock) {
	redisClient, mock := redismock.NewClientMock()
	return redisClient, mock
}

func TearDownTestRedis(mock redismock.ClientMock) error {
	return mock.ExpectationsWereMet()
}
