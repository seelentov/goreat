package repo

import (
	"goreat/internal/models/entities"
	"goreat/internal/models/queries"
)

type GetByResult struct {
	Entities []*entities.DBEntity
	Total    int64
	Exists   bool
	Error    error
}

type EntityRepository interface {
	GetByID(id uint) (*entities.DBEntity, error)
	ByQuery(query *queries.Query) GetByResult
	Create(topicName string, values map[string]interface{}) (*entities.DBEntity, error)
	UpdateByID(id uint, fTypes map[string]interface{}) error
	DeleteByID(id uint) error
	DeleteByQuery(query *queries.Query) error
}
