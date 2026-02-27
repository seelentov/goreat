package db

import "testing"

func TestNewInMemoryDB(t *testing.T) {
	db, err := NewInMemoryDB()
	if err != nil {
		t.Fatal(err)
	}
	if db == nil {
		t.Fatal("db is nil")
	}
}
