package db

import (
	"testing"

	"gorm.io/gorm"
)

var db *gorm.DB

func setup() {
	d, err := NewInMemoryDB()
	if err != nil {
		panic(err)
	}
	db = d
}

func TestClearDB(t *testing.T) {
	setup()

	if err := ClearDB(db); err != nil {
		t.Fatal(err)
	}
}

func TestSeedTestTopic(t *testing.T) {
	setup()

	if err := SeedTestTopic(db); err != nil {
		t.Fatal(err)
	}
}
