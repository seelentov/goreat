package repo

import (
	"goreat/internal/convension"
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
	err := e.db.Preload("Topic").Preload("Fields").First(&entity, id).Error
	if err != nil {
		return nil, err
	}

	fTypes, err := e.getScheme(entity.Topic.Name)
	if err != nil {
		return nil, err
	}

	if err := deserializeFields(&entity, fTypes); err != nil {
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
		v, err := convension.SerializeValue(v, fTypes[name].FieldType, fTypes[name].ContainerType)
		if err != nil {
			return nil, err
		}

		entity.Fields = append(entity.Fields, &entities.EntityField{
			Name:  name,
			Value: v,
		})
	}

	if err := deserializeFields(&entity, fTypes); err != nil {
		return nil, err
	}

	var topic entities.Topic
	if err := e.db.First(&topic, "name = ?", topicName).Error; err != nil {
		return nil, err
	}

	entity.TopicID = topic.ID

	return &entity, e.db.Create(&entity).Error
}

func (e *EntityRepositoryImpl) UpdateByID(id uint, values map[string]interface{}) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		var entity entities.Entity
		err := tx.Preload("Topic").Preload("Fields").First(&entity, id).Error
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

		for name, value := range values {
			index := fieldsIndexes[name]

			v, err := convension.SerializeValue(value, fTypes[name].FieldType, fTypes[name].ContainerType)
			if err != nil {
				return err
			}

			entity.Fields[index].Value = v
		}

		return tx.Save(&entity).Error
	})
}

func (e *EntityRepositoryImpl) DeleteByID(id uint) error {
	return e.db.Unscoped().Delete(&entities.Entity{}, id).Error
}

func (e *EntityRepositoryImpl) DeleteBy(topicName string, query *models.Query) error {
	panic("TODO: Implement")
}

func (e *EntityRepositoryImpl) getScheme(topicName string) (map[string]entities.FieldValueInfo, error) {
	var topic entities.Topic
	err := e.db.Preload("Fields").First(&topic, "name = ?", topicName).Error
	if err != nil {
		return nil, err
	}

	fTypes := make(map[string]entities.FieldValueInfo, len(topic.Fields))
	for _, f := range topic.Fields {
		fTypes[f.Name] = entities.FieldValueInfo{
			FieldType:     f.Type,
			ContainerType: f.ContainerType,
		}
	}

	return fTypes, nil
}

func deserializeFields(entity *entities.Entity, fTypes map[string]entities.FieldValueInfo) error {
	for _, f := range entity.Fields {
		v, err := convension.DeserializeValue(f.Value, fTypes[f.Name].FieldType, fTypes[f.Name].ContainerType)
		if err != nil {
			return err
		}
		f.ValueDecoded = v
	}
	return nil
}
