package entities

import (
	"gorm.io/gorm"
)

type DBEntity struct {
	gorm.Model

	TopicID uint
	Topic   *DBTopic `gorm:"foreignKey:TopicID"`

	Fields []*DBEntityField `gorm:"foreignKey:EntityID;constraint:OnDelete:CASCADE;"`
}

func NewDBEntity(values map[string]interface{}) (*DBEntity, error) {
	fields := make([]*DBEntityField, len(values))

	i := 0
	for name, value := range values {
		f, err := NewDBEntityField(name, value)
		if err != nil {
			return nil, err
		}

		fields[i] = f
		i++
	}

	return &DBEntity{
		Fields: fields,
	}, nil
}

func (e DBEntity) Flat() map[string]any {
	m := make(map[string]any, len(e.Fields))

	for _, f := range e.Fields {
		if f.Value == nil || len(f.Value) == 0 || f.Value[0] == nil {
			m[f.Name] = nil
			continue
		}

		if len(f.Value) == 1 {
			m[f.Name] = f.Value[0].GetValue()
			continue
		}

		valSlice := make([]any, len(f.Value))
		for i, value := range f.Value {
			valSlice[i] = value.GetValue()
		}
		m[f.Name] = valSlice
	}

	return m
}
