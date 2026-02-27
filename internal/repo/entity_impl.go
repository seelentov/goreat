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

func (e EntityRepositoryImpl) ByQuery(query *queries.Query) GetByResult {
	fTypes, err := e.getScheme(query.Topic)
	if err != nil {
		return GetByResult{
			Error: err,
		}
	}

	db, err := query.ToDB(e.db, fTypes)
	if err != nil {
		return GetByResult{
			Error: err,
		}
	}

	switch query.Type {
	case queries.QueryTypeCount:

	case queries.QueryTypeExists:
		var exists int
		err = db.Scan(&exists).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return GetByResult{
				Error: err,
			}
		}
		return GetByResult{
			Exists: exists != 0,
		}
	case queries.QueryTypeData:
		var ens []*entities.DBEntity
		if err := db.Find(&ens).Error; err != nil {
			return GetByResult{
				Error: err,
			}
		}
		return GetByResult{
			Entities: ens,
		}
	}

	return GetByResult{
		Error: entities.ErrType,
	}
}

func (e EntityRepositoryImpl) ExistsBy(query *queries.Query) ([]*entities.DBEntity, error) {
	fTypes, err := e.getScheme(query.Topic)
	if err != nil {
		return nil, err
	}

	db, err := query.ToDB(e.db, fTypes)
	if err != nil {
		return nil, err
	}

	var ens []*entities.DBEntity
	if err := db.Find(&ens).Error; err != nil {
		return nil, err
	}

	return ens, nil
}

func (e EntityRepositoryImpl) CountBy(query *queries.Query) ([]*entities.DBEntity, error) {
	fTypes, err := e.getScheme(query.Topic)
	if err != nil {
		return nil, err
	}

	db, err := query.ToDB(e.db, fTypes)
	if err != nil {
		return nil, err
	}

	var ens []*entities.DBEntity
	if err := db.Find(&ens).Error; err != nil {
		return nil, err
	}

	return ens, nil
}

func (e EntityRepositoryImpl) Create(topicName string, values map[string]interface{}) (*entities.DBEntity, error) {
	v, err := entities.NewDBEntity(values)
	if err != nil {
		return nil, err
	}

	topic, err := e.getTopic(topicName)
	if err != nil {
		return nil, err
	}

	v.Topic.ID = topic.ID

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

func (e EntityRepositoryImpl) DeleteByQuery(query *queries.Query) error {
	//TODO implement me
	panic("implement me")
}

func (e EntityRepositoryImpl) getTopic(topicName string) (*entities.DBTopic, error) {
	var topic *entities.DBTopic
	if err := e.db.First(&topic, "name = ?", topicName).Error; err != nil {
		return nil, err
	}
	return topic, nil
}

func (e EntityRepositoryImpl) getScheme(topicName string) (map[string]entities.FieldValueInfo, error) {
	var topic *entities.DBTopic
	if err := e.db.Preload("Fields").First(&topic, "name = ?", topicName).Error; err != nil {
		return nil, err
	}

	fTypes := make(map[string]entities.FieldValueInfo)
	for _, f := range topic.Fields {
		fTypes[f.Name] = entities.FieldValueInfo{
			FieldType:     f.Type,
			ContainerType: f.ContainerType,
		}
	}

	return fTypes, nil
}
