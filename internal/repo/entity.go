package repo

import (
	"goreat/internal/models"
	"goreat/internal/models/entities"
)

type EntityRepository interface {
	GetByID(id uint) (*entities.Entity, error)
	GetBy(topicName string, query *models.Query) ([]*entities.Entity, error)
	Create(topicName string, values map[string]interface{}) (*entities.Entity, error)
	UpdateByID(id uint, fTypes map[string]interface{}) error
	DeleteByID(id uint) error
	DeleteBy(topicName string, query *models.Query) error
}
