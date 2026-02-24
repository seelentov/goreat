package entities

import (
	"gorm.io/gorm"
)

type Entity struct {
	gorm.Model

	TopicID uint   `gorm:"index"`
	Topic   *Topic `gorm:"foreignKey:TopicID;references:ID"`

	Fields []*EntityField `gorm:"foreignKey:EntityID;constraint:OnDelete:CASCADE"`
}

func NewEntity(fTypes map[string][]byte) *Entity {
	fields := make([]*EntityField, 0, len(fTypes))

	for name, value := range fTypes {
		ef := &EntityField{
			Name:  name,
			Value: value,
		}
		if ef == nil {
			return nil
		}
		fields = append(fields, ef)
	}

	return &Entity{
		Fields: fields,
	}
}

type EntityField struct {
	EntityID uint    `gorm:"index"`
	Entity   *Entity `gorm:"foreignKey:EntityID;references:ID"`

	Name  string
	Value []byte

	ValueDecoded any `gorm:"-"`
}
