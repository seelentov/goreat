package repo

import "goreat/internal/models/entities"

type TopicRepository interface {
	GetAll() ([]*entities.DBTopic, error)
	GetByName(name string) (*entities.DBTopic, error)
	Create(name string, fTypes map[string][]byte) (*entities.DBTopic, error)
	Update(id uint, name string, fTypes map[string][]byte) (*entities.DBTopic, error)
}
