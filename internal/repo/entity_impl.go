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

func (e *EntityRepositoryImpl) GetByID(id uint) (*entities.DBEntity, error) {
	var entity entities.DBEntity
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

func (e *EntityRepositoryImpl) GetBy(topicName string, query *models.Query) ([]*entities.DBEntity, error) {
	panic("TODO: Implement")
}

func (e *EntityRepositoryImpl) Create(topicName string, values map[string]interface{}) (*entities.DBEntity, error) {
	fTypes, err := e.getScheme(topicName)
	if err != nil {
		return nil, err
	}

	var entity entities.DBEntity

	for name, v := range values {
		v, err := convension.SerializeValue(v, fTypes[name].FieldType, fTypes[name].ContainerType)
		if err != nil {
			return nil, err
		}

		entity.Fields = append(entity.Fields, &entities.DBEntityField{
			Name:  name,
			Value: v,
		})
	}

	if err := deserializeFields(&entity, fTypes); err != nil {
		return nil, err
	}

	var topic entities.DBTopic
	if err := e.db.First(&topic, &entities.DBTopic{Name: topicName}).Error; err != nil {
		return nil, err
	}

	entity.TopicID = topic.ID

	return &entity, e.db.Create(&entity).Error
}

func (e *EntityRepositoryImpl) UpdateByID(id uint, values map[string]interface{}) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		var entity entities.DBEntity
		if err := tx.Preload("Topic").Preload("Fields").First(&entity, id).Error; err != nil {
			return err
		}

		fTypes, err := e.getScheme(entity.Topic.Name)
		if err != nil {
			return err
		}

		for name, value := range values {
			v, err := convension.SerializeValue(value, fTypes[name].FieldType, fTypes[name].ContainerType)
			if err != nil {
				return err
			}

			if err := tx.Model(&entities.DBEntityField{}).Where(&entities.DBEntityField{EntityID: id, Name: name}).Updates(&entities.DBEntityField{Value: v}).Error; err != nil {
				return err
			}
		}

		return tx.Save(&entity).Error
	})
}

func (e *EntityRepositoryImpl) DeleteByID(id uint) error {
	return e.db.Unscoped().Delete(&entities.DBEntity{}, id).Error
}

func (e *EntityRepositoryImpl) DeleteBy(topicName string, query *models.Query) error {
	panic("TODO: Implement")
}

func (e *EntityRepositoryImpl) getScheme(topicName string) (map[string]entities.FieldValueInfo, error) {
	var topic entities.DBTopic
	err := e.db.Select("ID").Preload("Fields").First(&topic, &entities.DBTopic{Name: topicName}).Error
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

func deserializeFields(entity *entities.DBEntity, fTypes map[string]entities.FieldValueInfo) error {
	for _, f := range entity.Fields {
		v, err := convension.DeserializeValue(f.Value, fTypes[f.Name].FieldType, fTypes[f.Name].ContainerType)
		if err != nil {
			return err
		}
		f.ValueDecoded = v
	}
	return nil
}
