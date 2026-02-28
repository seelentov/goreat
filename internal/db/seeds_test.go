package db_test

import (
	"goreat/internal/db"
	"testing"

	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	d, err := db.NewSQLiteInMemoryDB()
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func TestClearDB(t *testing.T) {
	d := setupDB(t)

	if err := db.ClearDB(d); err != nil {
		t.Fatal(err)
	}
}

func TestSeedTestTopic(t *testing.T) {
	d := setupDB(t)

	if err := db.SeedTestTopic(d); err != nil {
		t.Fatal(err)
	}
}
