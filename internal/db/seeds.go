package db

import (
	"fmt"
	"goreat/internal/convension"
	"goreat/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

const TEST_TOPIC = "test"

func ClearDB(db *gorm.DB) error {
	for _, m := range dbModels {
		if err := db.Unscoped().Where("1 = 1").Delete(m).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedTestTopic(db *gorm.DB) error {
	fields := map[string]entities.FieldValueInfo{
		"string": {
			FieldType:     entities.FieldTypeString,
			ContainerType: entities.ContainerTypeSingle,
		},
		"int": {
			FieldType:     entities.FieldTypeInt,
			ContainerType: entities.ContainerTypeSingle,
		},
		"float": {
			FieldType:     entities.FieldTypeFloat,
			ContainerType: entities.ContainerTypeSingle,
		},
		"bool": {
			FieldType:     entities.FieldTypeBool,
			ContainerType: entities.ContainerTypeSingle,
		},
		"date": {
			FieldType:     entities.FieldTypeDateTime,
			ContainerType: entities.ContainerTypeSingle,
		},
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

		v := make(map[string][]byte, len(values))
		for fieldName, value := range values {
			serVal, err := convension.SerializeValue(value, fields[fieldName].FieldType, fields[fieldName].ContainerType)
			if err != nil {
				return err
			}
			v[fieldName] = serVal
		}

		entity := entities.NewEntity(v)
		entity.TopicID = topic.ID

		if err := db.Create(&entity).Error; err != nil {
			return err
		}
	}

	return nil
}
