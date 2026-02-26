package db

import (
	"goreat/internal/models/entities"
)

var entityDbModels = []interface{}{
	&entities.DBEntity{},
	&entities.DBEntityField{},
	&entities.DBEntityFieldValue{},

	&entities.DBTopic{},
	&entities.DBTopicField{},
}

var dbModels = append(
	entityDbModels,
	// &models.Page{},
	// &models.QueryData{},
)
