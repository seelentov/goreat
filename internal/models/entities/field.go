package models

import "gorm.io/gorm"

type FieldType int

const (
	FieldTypeString FieldType = iota
	FieldTypeInt
	FieldTypeFloat
	FieldTypeBool
	FieldTypeDateTime
)

type Field struct {
	gorm.Model

	Name string
	Type FieldType
}

func NewField(name string, fType FieldType) *Field {
	return &Field{
		Name: name,
		Type: fType,
	}
}
