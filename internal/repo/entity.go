package repo

import (
	"goreat/internal/models"
	"goreat/internal/models/entities"
)

type EntityRepository interface {
	GetByID(id uint) (*entities.Entity, error)
	GetBy(topicName string, query *models.Query) ([]*entities.Entity, error)
	Create(topicName string, fTypes map[string]entities.FieldValuePair) error
	UpdateByID(id uint, fTypes map[string]entities.FieldValuePair) error
	DeleteByID(id uint) error
	DeleteBy(topicName string, query *models.Query) error
}
