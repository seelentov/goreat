package entities

import "gorm.io/gorm"

type Topic struct {
	gorm.Model

	Name   string        `gorm:"uniqueIndex"`
	Fields []*TopicField `gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE"`

	Entities []*Entity `gorm:"many2many:topic_entities;"`
}

func NewTopic(name string, fields map[string]FieldValueInfo) *Topic {
	f := make([]*TopicField, 0, len(fields))
	for name, fieldType := range fields {
		f = append(f, &TopicField{
			Field: &Field{
				Name:          name,
				Type:          fieldType.FieldType,
				ContainerType: fieldType.ContainerType,
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

	TopicID uint   `gorm:"index"`
	Topic   *Topic `gorm:"foreignKey:TopicID;references:ID"`
}
