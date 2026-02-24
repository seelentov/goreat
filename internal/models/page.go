package models

import "gorm.io/gorm"

type DBPage struct {
	gorm.Model

	Path     string `gorm:"uniqueIndex"`
	Template string

	Queries []*Query `gorm:"type:json"`
}
