package repo

import "gorm.io/gorm"

type TopicRepositoryImpl struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepositoryImpl {
	return &TopicRepositoryImpl{
		db: db,
	}
}
