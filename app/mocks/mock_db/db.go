package mock_db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize GORM: %v", err)
	}

	return gormDB, mock
}

func TearDownTestDB(db *gorm.DB, mock sqlmock.Sqlmock) error {
	sqlDB, _ := db.DB()
	// TODO: rearrange the expect close clause to above the closing db
	err := sqlDB.Close()
	if err != nil {
		return err
	}
	mock.ExpectClose()

	return mock.ExpectationsWereMet()
}
