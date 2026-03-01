package repo_test

import (
	"errors"
	"fmt"
	"goreat/internal/db"
	"goreat/internal/models/queries"
	"goreat/internal/repo"
	"os"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
)

var database *gorm.DB

func setupEntityTest(t *testing.T) *repo.EntityRepository {
	t.Helper()

	tempDBFilePath := fmt.Sprintf("test_%d.db", time.Now().UnixNano())
	database, err := db.NewSQLiteFileDB(tempDBFilePath)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	if err := db.SeedTestTopic(database); err != nil {
		t.Fatalf("failed to seed db: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := database.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		_ = os.Remove(tempDBFilePath)
	})

	topicRepoImpl := repo.NewTopicRepository(database)
	return repo.NewEntityRepository(topicRepoImpl, database)
}

func TestEntityRepositoryImpl_GetByID(t *testing.T) {
	entityRepo := setupEntityTest(t)

	entity, err := entityRepo.GetByID(1)
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
	table := getValues(0)
	testEntity(flat, table, t)
}

func TestEntityRepositoryImpl_Create(t *testing.T) {
	entityRepo := setupEntityTest(t)
	i := 1

	entity, err := entityRepo.Create("test", getValues(i))
	if err != nil {
		t.Error(err)
	}

	newEntity, err := entityRepo.GetByID(entity.ID)
	if err != nil {
		t.Error(err)
	}

	if newEntity == nil {
		t.Error("entity is nil")
	}

	if newEntity.ID != entity.ID {
		t.Errorf("got ID %d, expected %d", entity.ID, entity.ID)
	}

	flat := newEntity.Flat()
	table := getValues(i)
	testEntity(flat, table, t)
}

func TestEntityRepositoryImpl_UpdateByID(t *testing.T) {
	entityRepo := setupEntityTest(t)
	i := 1

	entity, err := entityRepo.GetByID(1)
	if err != nil {
		t.Error(err)
	}

	if entity == nil {
		t.Error("entity is nil")
	}

	if err := entityRepo.UpdateByID(entity.ID, getValues(i)); err != nil {
		t.Error(err)
	}

	newEntity, err := entityRepo.GetByID(entity.ID)
	if err != nil {
		t.Error(err)
	}

	if newEntity == nil {
		t.Error("entity is nil")
	}

	if newEntity.ID != entity.ID {
		t.Errorf("got ID %d, expected %d", entity.ID, entity.ID)
	}

	flat := newEntity.Flat()
	table := getValues(i)
	testEntity(flat, table, t)
}

func TestEntityRepositoryImpl_DeleteByID(t *testing.T) {
	entityRepo := setupEntityTest(t)
	id := uint(1)

	if err := entityRepo.DeleteByID(id); err != nil {
		t.Error(err)
	}

	entity, err := entityRepo.GetByID(id)

	if err == nil {
		t.Error("error is nil")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Error(err)
	}

	if entity != nil {
		t.Error("entity is not nil")
	}
}

func TestEntityRepositoryImpl_ByQuery(t *testing.T) {
	entityRepo := setupEntityTest(t)

	limit := uint(10)
	q := queries.Query{
		Topic: "test",
		Filters: []queries.Filter{
			{
				Field: "string",
				Type:  queries.FilterTypeContains,
				Value: "1",
			},
			{
				Field: "int",
				Type:  queries.FilterTypeGreaterThan,
				Value: "10",
			},
		},
		Orders: []queries.Order{
			{
				Field:     "int",
				Direction: queries.OrderDirectionAsc,
			},
		},
		Limit: &limit,
		Type:  queries.QueryTypeData,
	}

	res := entityRepo.ByQuery(q)
	if res.Error != nil {
		t.Error(res.Error)
	}

	ens := res.Entities

	if len(ens) == 0 {
		t.Errorf("Expected %v entities, got %v", limit, 0)
	}

	if len(ens) > 10 {
		t.Errorf("Expected %v entities, got %v", limit, len(ens))
	}

	i := int64(0)
	for _, entity := range ens {
		en := entity.Flat()

		if i > en["int"].(int64) {
			t.Errorf("int field not ordered")
		} else {
			i = en["int"].(int64)
		}

		if en["int"].(int64) < 10 {
			t.Errorf("int field not filtered")
		}

		if !strings.Contains(en["string"].(string), "1") {
			t.Errorf("string field not filtered")
		}
	}
}

func testEntity(actual, expected map[string]interface{}, t *testing.T) {
	for key, value := range expected {
		if key == "date" {
			continue
		}

		if actual[key] != value {
			t.Errorf("got %v, expected %v", actual[key], value)
		}
	}

	if _, ok := expected["date"]; ok {
		expDate := time.Now().Add(time.Hour * time.Duration(0))
		if actual["date"].(time.Time).Equal(expDate) {
			t.Errorf("got %v, expected %v", actual["date"], expDate)
		}
	}
}

func getValues(i int) map[string]interface{} {
	m := make(map[string]interface{}, 4)

	m["string"] = fmt.Sprintf("string %v", i)
	m["int"] = int64(i)
	m["float"] = float64(i) / 1000.0
	m["bool"] = i%2 == 0
	m["date"] = time.Now().Add(time.Hour * time.Duration(i))

	return m
}
