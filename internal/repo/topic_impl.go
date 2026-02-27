package repo

import (
	"goreat/internal/models/entities"

	"gorm.io/gorm"
)

type TopicRepositoryImpl struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepositoryImpl {
	return &TopicRepositoryImpl{
		db: db,
	}
}

func (t *TopicRepositoryImpl) GetAll() ([]*entities.DBTopic, error) {
	panic("TODO: Implement")
}

func (t *TopicRepositoryImpl) GetByName(name string) (*entities.DBTopic, error) {
	panic("TODO: Implement")
}

func (t *TopicRepositoryImpl) Create(name string, fTypes map[string][]byte) (*entities.DBTopic, error) {
	panic("TODO: Implement")
}

func (t *TopicRepositoryImpl) Update(id uint, name string, fTypes map[string][]byte) (*entities.DBTopic, error) {
	panic("TODO: Implement")
}
