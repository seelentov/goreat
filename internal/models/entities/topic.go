package models

import "gorm.io/gorm"

type Topic struct {
	gorm.Model

	Name   string
	Fields []*TopicField

	Entities []*Entity
}

func NewTopic(name string, fields map[string]FieldType) *Topic {
	f := make([]*TopicField, 0, len(fields))
	for name, fieldType := range fields {
		f = append(f, &TopicField{
			Field: &Field{
				Name: name,
				Type: fieldType,
			},
		})
	}

	return &Topic{
		Name:   name,
		Fields: f,
	}
}

type TopicField struct {
	*Field

	TopicID uint
	Topic   *Topic
}
