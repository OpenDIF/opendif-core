package testutils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupMockDB creates a mock database for unit testing using go-sqlmock.
// This should be used in unit tests where you want to mock database interactions.
//
// Usage:
//
//	db, mock, cleanup := testutils.SetupMockDB(t)
//	defer cleanup()
//	// Configure mock expectations
//	mock.ExpectQuery(...)
//	// Use db in your tests
func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		db.Close()
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Logf("warning: failed to close mock db: %v", err)
		}
	}

	return gormDB, mock, cleanup
}

