package repo

import "goreat/internal/models/entities"

type TopicRepository interface {
	GetAll() ([]*entities.DBTopic, error)
	GetByName(name string) (*entities.DBTopic, error)
	CreateOrUpdateByName(name string, fTypes map[string]entities.FieldInfo) (*entities.DBTopic, error)
	UpdateByName(name string, fTypes map[string]entities.FieldInfo) error
	DeleteByName(name string) error

	GetScheme(topicName string) (map[string]entities.FieldInfo, error)
}
