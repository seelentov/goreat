package db

import (
	"goreat/internal/models/entities"
)

var entityDbModels = []interface{}{
	&entities.Entity{},
	&entities.EntityField{},
	&entities.Topic{},
	&entities.TopicField{},
}

var dbModels = append(
	entityDbModels,
	// &models.Page{},
	// &models.QueryData{},
)
