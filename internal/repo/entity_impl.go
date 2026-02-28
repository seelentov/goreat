package repo

import (
	"goreat/internal/models/entities"
	"goreat/internal/models/queries"

	"gorm.io/gorm"
)

type EntityRepositoryImpl struct {
	topicRepo TopicRepository
	db        *gorm.DB
}

func NewEntityRepositoryImpl(topicRepo TopicRepository, db *gorm.DB) EntityRepository {
	return &EntityRepositoryImpl{
		topicRepo: topicRepo,
		db:        db,
	}
}

func (e EntityRepositoryImpl) GetByID(id uint) (*entities.DBEntity, error) {
	var entity entities.DBEntity
	if err := e.db.Preload("Fields.Value").Preload("Topic").First(&entity, id).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (e EntityRepositoryImpl) ByQuery(query queries.Query) GetByResult {
	fTypes, err := e.topicRepo.GetScheme(query.Topic)
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
	case queries.QueryTypeData:
		db = db.Select("db_entities.*")
	case queries.QueryTypeCount:
		db = db.Select("count(*)")
	case queries.QueryTypeExists:
		db = db.Select("1")
	default:
		return GetByResult{
			Error: entities.ErrType,
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

func (e EntityRepositoryImpl) Create(topicName string, values map[string]interface{}) (*entities.DBEntity, error) {
	v, err := entities.NewDBEntity(values)
	if err != nil {
		return nil, err
	}

	topic, err := e.topicRepo.GetByName(topicName)
	if err != nil {
		return nil, err
	}

	v.Topic = topic

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

	scheme, err := e.topicRepo.GetScheme(entity.Topic.Name)
	if err != nil {
		return err
	}

	for fName, _ := range fTypes {
		if _, ok := scheme[fName]; !ok {
			return ErrField
		}
	}

	fieldIndexes := make(map[string]int, len(entity.Fields))
	for i, f := range entity.Fields {
		fieldIndexes[f.Name] = i
	}

	return e.db.Transaction(func(tx *gorm.DB) error {
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
	})
}

func (e EntityRepositoryImpl) DeleteByID(id uint) error {
	return e.db.Delete(&entities.DBEntity{}, id).Error
}
