package db

import (
	"fmt"
	"goreat/internal/models/entities"
	"time"

	"gorm.io/gorm"
)

func ClearDB(db *gorm.DB) error {
	for _, m := range dbModels {
		if err := db.Unscoped().Where("1 = 1").Delete(m).Error; err != nil {
			return err
		}
	}
	return nil
}

var TestTopicFields = map[string]entities.FieldValueInfo{
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
	"string_list": {
		FieldType:     entities.FieldTypeString,
		ContainerType: entities.ContainerTypeList,
	},
}

func SeedTestTopic(db *gorm.DB) error {
	topic := entities.NewTopic("test", TestTopicFields)

	if err := db.Create(topic).Error; err != nil {
		return err
	}

	for i := range 100 {
		values := map[string]interface{}{
			"string":      fmt.Sprintf("string %v", i),
			"int":         i,
			"float":       float64(i) / 1000.0,
			"bool":        i%2 == 0,
			"date":        time.Now().Add(time.Hour * time.Duration(i)),
			"string_list": []string{fmt.Sprintf("string %v", i*10), fmt.Sprintf("string %v", i*20), fmt.Sprintf("string %v", i*30)},
		}

		entity, err := entities.NewDBEntity(values)
		if err != nil {
			return err
		}

		if err := db.Create(&entity).Error; err != nil {
			return err
		}
	}

	return nil
}
