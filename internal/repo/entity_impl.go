package repo

import (
	"goreat/internal/models"
	"goreat/internal/models/entities"

	"gorm.io/gorm"
)

type EntityRepositoryImpl struct {
	db *gorm.DB
}

func (e *EntityRepositoryImpl) GetByID(id uint) (*entities.Entity, error) {
	var entity entities.Entity
	err := e.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (e *EntityRepositoryImpl) GetBy(topicName string, query *models.Query) ([]*entities.Entity, error) {
	panic("TODO: Implement")
}

func (e *EntityRepositoryImpl) Create(topicName string, fTypes map[string]entities.FieldValuePair) error {
	var entity entities.Entity

	for name, v := range fTypes {
		entity.Fields = append(entity.Fields, entities.NewEntityField(name, v.FieldType, v.Value))
	}

	return e.db.Save(entity).Error
}

func (e *EntityRepositoryImpl) UpdateByID(id uint, fTypes map[string]any) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		var entity entities.Entity
		err := tx.First(&entity, id).Error
		if err != nil {
			return err
		}

		fieldsIndexes := make(map[string]int, len(entity.Fields))
		for i, f := range entity.Fields {
			fieldsIndexes[f.Name] = i
		}

		for name, value := range fTypes {
			index := fieldsIndexes[name]
			entity.Fields[index].SetValue(value)
		}

		return tx.Save(entity).Error
	})
}

func (e *EntityRepositoryImpl) DeleteByID(id uint) error {
	return e.db.Delete(&entities.Entity{}, id).Error
}

func (e *EntityRepositoryImpl) DeleteBy(topicName string, query *models.Query) error {
	panic("TODO: Implement")
}
