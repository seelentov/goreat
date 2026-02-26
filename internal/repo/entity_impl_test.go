package repo

import (
	"goreat/internal/db"
	"testing"
	"time"

	"gorm.io/gorm"
)

var entityRepoImpl EntityRepository
var database *gorm.DB

func setup(t *testing.T) {
	d, err := db.NewInMemoryDB()
	if err != nil {
		t.Error(err)
	}

	database = d

	err = db.SeedTestTopic(database)
	if err != nil {
		t.Error(err)
	}

	entityRepoImpl = NewEntityRepositoryImpl(database)
}

func teardown(t *testing.T) {
	if err := db.ClearDB(database); err != nil {
		t.Error(err)
	}
}

func TestEntityRepositoryImpl_GetByID(t *testing.T) {
	setup(t)

	entity, err := entityRepoImpl.GetByID(1)
	if err != nil {
		t.Error(err)
	}

	if entity == nil {
		t.Error("entity is nil")
	}

	if entity.ID != 1 {
		t.Errorf("got ID %d, expected 1", entity.ID)
	}

	flat := entity.Flat()

	table := []struct {
		key   string
		value any
	}{
		{
			key:   "string",
			value: "string 0",
		},
		{
			key:   "int",
			value: int64(0),
		},
		{
			key:   "float",
			value: float64(0) / 1000.0,
		},
		{
			key:   "bool",
			value: true,
		},
	}

	for _, ta := range table {
		if flat[ta.key] != ta.value {
			t.Errorf("got %v, expected %v", flat[ta.key], ta.value)
		}
	}

	expDate := time.Now().Add(time.Hour * time.Duration(0))
	if flat["date"].(time.Time).Equal(expDate) {
		t.Errorf("got %v, expected %v", flat["date"], expDate)
	}

	teardown(t)
}
