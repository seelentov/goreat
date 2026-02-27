package repo

import (
	"goreat/internal/models/entities"
	"goreat/internal/models/queries"

	"gorm.io/gorm"
)

type EntityRepositoryImpl struct {
	db *gorm.DB
}

func NewEntityRepositoryImpl(db *gorm.DB) EntityRepository {
	return &EntityRepositoryImpl{
		db: db,
	}
}

func (e EntityRepositoryImpl) GetByID(id uint) (*entities.DBEntity, error) {
	var entity entities.DBEntity
	if err := e.db.Preload("Fields.Value").First(&entity, id).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (e EntityRepositoryImpl) GetBy(topicName string, query *queries.Query) ([]*entities.DBEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (e EntityRepositoryImpl) Create(topicName string, values map[string]interface{}) (*entities.DBEntity, error) {
	v, err := entities.NewDBEntity(values)
	if err != nil {
		return nil, err
	}
	if err := e.db.Create(v).Error; err != nil {
		return nil, err
	}
	return v, nil
}

func (e EntityRepositoryImpl) UpdateByID(id uint, fTypes map[string]interface{}) error {
	entity, err := e.GetByID(id)
	if err != nil {
		return err
	}

	fieldIndexes := make(map[string]int, len(entity.Fields))
	for i, f := range entity.Fields {
		fieldIndexes[f.Name] = i
	}

	for key, value := range fTypes {
		idx, ok := fieldIndexes[key]
		if !ok {
			continue
		}

		field := &entity.Fields[idx]

		if err := (*field).SetValue(value); err != nil {
			return err
		}

		err = e.db.Model(field).Association("Value").Replace((*field).Value)
		if err != nil {
			return err
		}
	}

	return e.db.Save(entity).Error
}

func (e EntityRepositoryImpl) DeleteByID(id uint) error {
	return e.db.Delete(&entities.DBEntity{}, id).Error
}

func (e EntityRepositoryImpl) DeleteBy(topicName string, query *queries.Query) error {
	//TODO implement me
	panic("implement me")
}
