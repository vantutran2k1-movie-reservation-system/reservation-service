package transaction

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ExecuteInRedisTransaction(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
	err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		return fn(tx)
	})

	return err
}
