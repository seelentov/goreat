package db

import (
	"fmt"
	"goreat/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

const TEST_TOPIC = "test"

func SeedTestTopic(db *gorm.DB) error {
	fields := map[string]entities.FieldType{
		"string": entities.FieldTypeString,
		"int":    entities.FieldTypeInt,
		"float":  entities.FieldTypeFloat,
		"bool":   entities.FieldTypeBool,
		"date":   entities.FieldTypeDateTime,
	}

	topic := entities.NewTopic("test", fields)

	if err := db.Create(topic).Error; err != nil {
		return err
	}

	for i := range 10 {
		values := map[string]interface{}{
			"string": fmt.Sprintf("string %v", i),
			"int":    i,
			"float":  float64(i) / 1000.0,
			"bool":   i%2 == 0,
			"date":   time.Now().Add(time.Hour * time.Duration(i)),
		}

		v := make(map[string]entities.FieldValuePair, len(values))
		for fieldName, value := range values {
			v[fieldName] = entities.FieldValuePair{
				FieldType: fields[fieldName],
				Value:     value,
			}
		}

		entity := entities.NewEntity(v)

		if err := db.Create(&entity).Error; err != nil {
			return err
		}
	}

	return nil
}
