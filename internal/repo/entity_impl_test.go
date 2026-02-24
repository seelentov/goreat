package repo

import (
	"fmt"
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

	entityRepoImpl = NewEntityRepository(database)
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
		t.Error("entity id is not 1")
	}

	teardown(t)
}

func TestEntityRepositoryImpl_Create(t *testing.T) {
	setup(t)

	i := 99
	values := map[string]interface{}{
		"string": fmt.Sprintf("string %v", i),
		"int":    int64(i),
		"float":  float64(i) / 1000.0,
		"bool":   i%2 == 0,
		"date":   time.Now().Add(time.Hour * time.Duration(i)),
	}

	entity, err := entityRepoImpl.Create(db.TEST_TOPIC, values)
	if err != nil {
		t.Error(err)
	}

	if entity == nil {
		t.Error("new entity var is empty")
	}

	newID := entity.ID

	entity, err = entityRepoImpl.GetByID(newID)
	if err != nil {
		t.Error(err)
	}
	if entity == nil {
		t.Error("entity is nil")
	}
	if entity.ID != newID {
		t.Errorf("entity id is not %v\n", newID)
	}

	for _, f := range entity.Fields {
		if f.Name == "date" {
			if !f.ValueDecoded.(time.Time).Equal(values[f.Name].(time.Time)) {
				t.Errorf("field %v is %v, expected %v", f.Name, f.ValueDecoded, values[f.Name])
			}

			continue
		}

		if f.ValueDecoded != values[f.Name] {
			t.Errorf("field %v is %v, expected %v", f.Name, f.ValueDecoded, values[f.Name])
		}
	}

	teardown(t)
}

func TestEntityRepositoryImpl_UpdateByID(t *testing.T) {
	setup(t)

	i := 99

	values := map[string]interface{}{
		"string": fmt.Sprintf("string %v", i),
		"int":    int64(i),
		"float":  float64(i) / 1000.0,
		"bool":   i%2 == 0,
		"date":   time.Now().Add(time.Hour * time.Duration(i)),
	}

	if err := entityRepoImpl.UpdateByID(1, values); err != nil {
		t.Error(err)
	}

	entity, err := entityRepoImpl.GetByID(1)
	if err != nil {
		t.Error(err)
	}

	if entity == nil {
		t.Error("entity is nil")
	}

	if entity.ID != 1 {
		t.Errorf("entity id is not 1\n")
	}

	for _, f := range entity.Fields {
		if f.Name == "date" {
			if !f.ValueDecoded.(time.Time).Equal(values[f.Name].(time.Time)) {
				t.Errorf("field %v is %v, expected %v", f.Name, f.ValueDecoded, values[f.Name])
			}

			continue
		}

		if f.ValueDecoded != values[f.Name] {
			t.Errorf("field %v is %v, expected %v", f.Name, f.ValueDecoded, values[f.Name])
		}
	}

	teardown(t)
}

func TestEntityRepositoryImpl_DeleteByID(t *testing.T) {
	setup(t)

	if err := entityRepoImpl.DeleteByID(1); err != nil {
		t.Error(err)
	}

	entity, err := entityRepoImpl.GetByID(1)
	if err != gorm.ErrRecordNotFound {
		t.Error(err)
	}
	if entity != nil {
		t.Error("entity is not nil")
	}

	teardown(t)
}

func TestEntityRepositoryImpl_GetBy(t *testing.T) {
	setup(t)

	teardown(t)
}

func TestEntityRepositoryImpl_DeleteBy(t *testing.T) {
	setup(t)

	teardown(t)
}
