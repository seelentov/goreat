package repo

import (
	"goreat/internal/models/entities"

	"gorm.io/gorm"
)

type TopicRepositoryImpl struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return &TopicRepositoryImpl{
		db: db,
	}
}

func (t TopicRepositoryImpl) GetAll() ([]*entities.DBTopic, error) {
	var topics []*entities.DBTopic
	if err := t.db.Preload("Fields").Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func (t TopicRepositoryImpl) GetByName(name string) (*entities.DBTopic, error) {
	var topic *entities.DBTopic
	if err := t.db.Preload("Fields").First(&topic, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return topic, nil
}

func (t TopicRepositoryImpl) CreateOrUpdateByName(name string, fTypes map[string]entities.FieldInfo) (*entities.DBTopic, error) {
	//TODO implement me
	panic("implement me")
}

func (t TopicRepositoryImpl) UpdateByName(name string, fTypes map[string]entities.FieldInfo) error {
	//TODO implement me
	panic("implement me")
}

func (t TopicRepositoryImpl) DeleteByName(name string) error {
	return t.db.Delete(&entities.DBTopic{}, "name = ?", name).Error
}

func (t TopicRepositoryImpl) GetScheme(topicName string) (map[string]entities.FieldInfo, error) {
	var topic *entities.DBTopic
	if err := t.db.Preload("Fields").First(&topic, "name = ?", topicName).Error; err != nil {
		return nil, err
	}

	fTypes := make(map[string]entities.FieldInfo)
	for _, f := range topic.Fields {
		fTypes[f.Name] = entities.FieldInfo{
			FieldType:     f.FieldType,
			ContainerType: f.ContainerType,
		}
	}

	return fTypes, nil
}
