package entities

import "gorm.io/gorm"

type FieldType int

const (
	FieldTypeString FieldType = iota
	FieldTypeInt
	FieldTypeFloat
	FieldTypeBool
	FieldTypeDateTime

	FieldTypeSource
)

type ContainerType int

const (
	ContainerTypeSingle ContainerType = iota
	ContainerTypeList
)

type Field struct {
	gorm.Model

	Name          string
	Type          FieldType
	ContainerType ContainerType
}

func NewField(name string, fType FieldType, cType ContainerType) *Field {
	return &Field{
		Name:          name,
		Type:          fType,
		ContainerType: cType,
	}
}

type FieldValueInfo struct {
	FieldType     FieldType
	ContainerType ContainerType
}

type FieldValueData struct {
	*FieldValueInfo

	Value any
}
