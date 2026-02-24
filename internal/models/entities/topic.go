package entities

import "gorm.io/gorm"

type DBTopic struct {
	gorm.Model

	Name   string          `gorm:"uniqueIndex"`
	Fields []*DBTopicField `gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE"`

	Entities []*DBEntity `gorm:"many2many:topic_entities;"`
}

func NewTopic(name string, fields map[string]FieldValueInfo) *DBTopic {
	f := make([]*DBTopicField, 0, len(fields))
	for name, fieldType := range fields {
		f = append(f, &DBTopicField{
			Field: &Field{
				Name:          name,
				Type:          fieldType.FieldType,
				ContainerType: fieldType.ContainerType,
			},
		})
	}

	return &DBTopic{
		Name:   name,
		Fields: f,
	}
}

type DBTopicField struct {
	*Field

	TopicID uint     `gorm:"index"`
	Topic   *DBTopic `gorm:"foreignKey:TopicID;references:ID"`
}
