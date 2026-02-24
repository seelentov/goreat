package models

import "gorm.io/gorm"

type Page struct {
	gorm.Model

	Path     string
	Template string
}
