package entities

import "gorm.io/gorm"

type DBTopic struct {
	gorm.Model

	Name   string          `gorm:"uniqueIndex"`
	Fields []*DBTopicField `gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE"`

	Entities []*DBEntity `gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE"`
}

func NewTopic(name string, fields map[string]FieldInfo) *DBTopic {
	fis := make([]*DBTopicField, len(fields))
	i := 0
	for name, fieldType := range fields {
		f := &DBTopicField{
			DefinitionField: DefinitionField{
				Field: Field{
					Name: name,
				},
				FieldInfo: FieldInfo{
					FieldType:     fieldType.FieldType,
					ContainerType: fieldType.ContainerType,
				},
			},
		}
		fis[i] = f
		i++
	}

	return &DBTopic{
		Name:   name,
		Fields: fis,
	}
}
