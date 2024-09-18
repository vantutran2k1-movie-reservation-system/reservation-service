package errors

import (
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsRedisKeyNotFoundError(err error) bool {
	return errors.Is(err, redis.Nil)
}
