package repo

import (
	"goreat/internal/models"
	"goreat/internal/models/entities"

	"gorm.io/gorm"
)

type EntityRepositoryImpl struct {
	db *gorm.DB
}

func NewEntityRepository(db *gorm.DB) *EntityRepositoryImpl {
	return &EntityRepositoryImpl{
		db: db,
	}
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

func (e *EntityRepositoryImpl) Create(topicName string, values map[string]interface{}) (*entities.Entity, error) {
	fTypes, err := e.getScheme(topicName)
	if err != nil {
		return nil, err
	}

	var entity entities.Entity

	for name, v := range values {
		entity.Fields = append(entity.Fields, entities.NewEntityField(name, fTypes[name], v))
	}

	return &entity, e.db.Create(&entity).Error
}

func (e *EntityRepositoryImpl) UpdateByID(id uint, values map[string]interface{}) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		var entity entities.Entity
		err := tx.Preload("Topic").First(&entity, id).Error
		if err != nil {
			return err
		}

		fTypes, err := e.getScheme(entity.Topic.Name)
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

func (e *EntityRepositoryImpl) getScheme(topicName string) (map[string]entities.FieldType, error) {
	var topic entities.Topic
	err := e.db.Preload("Fields").First(&topic, "name = ?", topicName).Error
	if err != nil {
		return nil, err
	}

	fTypes := make(map[string]entities.FieldType, len(topic.Fields))
	for _, f := range topic.Fields {
		fTypes[f.Name] = f.Type
	}

	return fTypes, nil
}
