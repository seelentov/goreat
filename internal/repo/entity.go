package repo

import (
	"goreat/internal/models/entities"
	"goreat/internal/models/queries"
)

type EntityRepository interface {
	GetByID(id uint) (*entities.DBEntity, error)
	GetBy(topicName string, query *queries.Query) ([]*entities.DBEntity, error)
	Create(topicName string, values map[string]interface{}) (*entities.DBEntity, error)
	UpdateByID(id uint, fTypes map[string]interface{}) error
	DeleteByID(id uint) error
	DeleteBy(topicName string, query *queries.Query) error
}
