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

func TearDownTestDB(db *gorm.DB, mock sqlmock.Sqlmock) {
	sqlDB, _ := db.DB()
	sqlDB.Close()
	mock.ExpectClose()
}
