package db_test

import (
	"goreat/internal/db"
	"os"
	"testing"
)

func TestNewSQLiteInMemoryDB(t *testing.T) {
	d, err := db.NewSQLiteInMemoryDB()
	if err != nil {
		t.Fatal(err)
	}
	if d == nil {
		t.Fatal("db is nil")
	}
}

func TestNewSQLiteFileDB(t *testing.T) {
	fileName := "test.db"
	d, err := db.NewSQLiteFileDB(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if d == nil {
		t.Fatal("db is nil")
	}

	t.Cleanup(func() {
		sqlDB, err := d.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		_ = os.Remove(fileName)
	})
}
