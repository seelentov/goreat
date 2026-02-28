package repo

import (
	"goreat/internal/models/entities"

	"gorm.io/gorm"
)

type TopicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{
		db: db,
	}
}

func (t TopicRepository) GetAll() ([]*entities.DBTopic, error) {
	var topics []*entities.DBTopic
	if err := t.db.Preload("Fields").Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func (t TopicRepository) GetByName(name string) (*entities.DBTopic, error) {
	var topic *entities.DBTopic
	if err := t.db.Preload("Fields").First(&topic, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return topic, nil
}

func (t TopicRepository) CreateOrUpdateByName(name string, fTypes map[string]entities.FieldInfo) (*entities.DBTopic, error) {
	topic, err := t.GetByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			topic = entities.NewTopic(name, fTypes)
			if err := t.db.Create(topic).Error; err != nil {
				return nil, err
			}
			return topic, nil
		}
		return nil, err
	}

	if err := t.UpdateByName(name, fTypes); err != nil {
		return nil, err
	}

	return t.GetByName(name)
}

func (t TopicRepository) UpdateByName(name string, fTypes map[string]entities.FieldInfo) error {
	topic, err := t.GetByName(name)
	if err != nil {
		return err
	}

	return t.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("topic_id = ?", topic.ID).Delete(&entities.DBTopicField{}).Error; err != nil {
			return err
		}

		for fName, fInfo := range fTypes {
			newField := entities.NewDBTopicField(fName, fInfo.FieldType, fInfo.ContainerType)
			newField.TopicID = topic.ID
			if err := tx.Create(newField).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (t TopicRepository) DeleteByName(name string) error {
	return t.db.Delete(&entities.DBTopic{}, "name = ?", name).Error
}

func (t TopicRepository) GetScheme(topicName string) (map[string]entities.FieldInfo, error) {
	topic, err := t.GetByName(topicName)
	if err != nil {
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
