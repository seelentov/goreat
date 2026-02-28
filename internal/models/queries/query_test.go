package queries

import (
	"fmt"
	"goreat/internal/db"
	"goreat/internal/models/entities"
	"os"
	"strings"
	"testing"

	"gorm.io/gorm"
)

var database *gorm.DB

var tempDBFilePath string

func setup() {
	tempDBFilePath = fmt.Sprintf("test_%v.db", "query")
	d, err := db.NewFileDB(tempDBFilePath)
	if err != nil {
		panic(err)
	}

	database = d

	err = db.SeedTestTopic(database)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	defer os.Remove(tempDBFilePath)

	sqlDB, err := database.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
}

func TestQuery_ToDB(t *testing.T) {
	setup()
	defer teardown()

	limit := uint(10)
	q := Query{
		Topic: "test",
		Filters: []*Filter{
			{
				Field: "string",
				Type:  FilterTypeContains,
				Value: "1",
			},
			{
				Field: "int",
				Type:  FilterTypeGreaterThan,
				Value: "10",
			},
		},
		Orders: []*Order{
			{
				Field:     "int",
				Direction: OrderDirectionAsc,
			},
		},
		Limit: &limit,
		Type:  QueryTypeData,
	}

	toDB, err := q.ToDB(database, db.TestTopicFields)
	if err != nil {
		t.Error(err)
	}

	toDB = toDB.Select("db_entities.*")

	var ens []*entities.DBEntity
	if err := toDB.Find(&ens).Error; err != nil {
		t.Error(err)
	}

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

func TestQuery_ToDB_Count(t *testing.T) {
	setup()
	defer teardown()

	q := Query{
		Topic: "test",
		Filters: []*Filter{
			{
				Field: "int",
				Type:  FilterTypeGreaterThan,
				Value: "10",
			},
		},
		Type: QueryTypeCount,
	}

	toDB, err := q.ToDB(database, db.TestTopicFields)
	if err != nil {
		t.Error(err)
	}

	toDB = toDB.Select("count(*)")

	var count int
	if err := toDB.Model(&entities.DBEntity{}).Scan(&count).Error; err != nil {
		t.Error(err)
	}

	if count == 0 {
		t.Errorf("Expected count > 0, got %v", count)
	}
}

func TestQuery_ToDB_Exists(t *testing.T) {
	setup()
	defer teardown()

	q := Query{
		Topic: "test",
		Filters: []*Filter{
			{
				Field: "string",
				Type:  FilterTypeEquals,
				Value: "not_existing_value_123",
			},
		},
		Type: QueryTypeExists,
	}

	toDB, err := q.ToDB(database, db.TestTopicFields)
	if err != nil {
		t.Error(err)
	}

	toDB = toDB.Select("1")

	var exists int
	err = toDB.Scan(&exists).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		t.Error(err)
	}

	if exists != 0 {
		t.Errorf("Expected not exists (0), got %v", exists)
	}

	q2 := Query{
		Topic: "test",
		Type:  QueryTypeExists,
	}

	toDB2, err := q2.ToDB(database, db.TestTopicFields)
	if err != nil {
		t.Error(err)
	}

	toDB2 = toDB.Select("1")

	err = toDB2.Scan(&exists).Error
	if err != nil {
		t.Error(err)
	}

	if exists != 1 {
		t.Errorf("Expected exists (1), got %v", exists)
	}
}

func TestQuery_ToDB_OffsetAndDescOrder(t *testing.T) {
	setup()
	defer teardown()

	limit := uint(5)
	offset := uint(2)
	q := Query{
		Topic: "test",
		Orders: []*Order{
			{
				Field:     "int",
				Direction: OrderDirectionDesc,
			},
		},
		Limit:  &limit,
		Offset: &offset,
		Type:   QueryTypeData,
	}

	toDB, err := q.ToDB(database, db.TestTopicFields)
	if err != nil {
		t.Error(err)
	}

	var ens []entities.DBEntity
	if err := toDB.Find(&ens).Error; err != nil {
		t.Error(err)
	}

	if len(ens) > int(limit) {
		t.Errorf("Expected max %v entities, got %v", limit, len(ens))
	}

	if len(ens) > 1 {
		en1 := ens[0].Flat()
		en2 := ens[1].Flat()

		if en1["int"].(int64) < en2["int"].(int64) {
			t.Errorf("Entities are not ordered in DESC direction")
		}
	}
}
