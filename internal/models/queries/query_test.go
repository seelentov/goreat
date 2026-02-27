package queries

import (
	"fmt"
	"goreat/internal/db"
	"goreat/internal/models/entities"
	"os"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
)

var database *gorm.DB

var tempDBFilePath string

func setup() {
	tempDBFilePath = fmt.Sprintf("test_%d.db", time.Now().UnixNano())
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
		Type:  QueryTypeCount,
	}

	toDB, err := q.ToDB(database)
	if err != nil {
		t.Error(err)
	}

	var ens []entities.DBEntity
	if err := toDB.Find(&ens).Error; err != nil {
		t.Error(err)
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
