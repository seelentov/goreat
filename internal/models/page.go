package models

import (
	"goreat/internal/models/queries"

	"gorm.io/gorm"
)

type DBPage struct {
	gorm.Model

	Path     string
	Template string

	Queries []queries.Query `gorm:"type:json"`
}
