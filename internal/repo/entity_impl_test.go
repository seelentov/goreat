package repo

import (
	"fmt"
	"goreat/internal/db"
	"testing"
	"time"
)

var entityRepoImpl EntityRepository

func setup() {
	d, err := db.NewInMemoryDB()
	if err != nil {
		panic(err)
	}

	err = db.SeedTestTopic(d)
	if err != nil {
		panic(err)
	}

	entityRepoImpl = NewEntityRepository(d)
}

func TestEntityRepositoryImpl_GetByID(t *testing.T) {
	setup()

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
}

func TestEntityRepositoryImpl_Create(t *testing.T) {
	setup()

	i := 99
	values := map[string]interface{}{
		"string": fmt.Sprintf("string %v", i),
		"int":    i,
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
		v, err := f.GetValue()
		if err != nil {
			t.Error(err)
		}

		if v != values[f.Name] {
			t.Errorf("field %v is not %v", f.Name, values[f.Name])
		}
	}
}
