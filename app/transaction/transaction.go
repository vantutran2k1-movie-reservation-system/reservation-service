package transaction

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TransactionManager interface {
	ExecuteInTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error
	ExecuteInRedisTransaction(rdb *redis.Client, fn func(tx *redis.Tx) error) error
}

func NewTransactionManager() TransactionManager {
	return &transactionManager{}
}

type transactionManager struct{}

func (m *transactionManager) ExecuteInTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (m *transactionManager) ExecuteInRedisTransaction(rdb *redis.Client, fn func(tx *redis.Tx) error) error {
	err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		return fn(tx)
	})

	return err
}
