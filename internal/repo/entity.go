package repo

import (
	"goreat/internal/models"
	"goreat/internal/models/entities"
)

type EntityRepository interface {
	GetByID(id uint) (*entities.DBEntity, error)
	GetBy(topicName string, query *models.Query) ([]*entities.DBEntity, error)
	Create(topicName string, values map[string]interface{}) (*entities.DBEntity, error)
	UpdateByID(id uint, fTypes map[string]interface{}) error
	DeleteByID(id uint) error
	DeleteBy(topicName string, query *models.Query) error
}
