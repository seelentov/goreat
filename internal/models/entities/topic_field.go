package entities

import "gorm.io/gorm"

type DBTopicField struct {
	gorm.Model

	DefinitionField

	TopicID uint
	Topic   *DBTopic `gorm:"foreignKey:TopicID"`
}

func NewDBTopicField(name string, fType FieldType, cType ContainerType) *DBTopicField {
	f := &DBTopicField{
		DefinitionField: *NewDefinitionField(name, fType, cType),
	}

	return f
}
