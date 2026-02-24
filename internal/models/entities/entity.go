package entities

import (
	"gorm.io/gorm"
)

type DBEntity struct {
	gorm.Model

	TopicID uint     `gorm:"index"`
	Topic   *DBTopic `gorm:"foreignKey:TopicID;references:ID"`

	Fields []*DBEntityField `gorm:"foreignKey:EntityID;constraint:OnDelete:CASCADE"`
}

func NewEntity(fTypes map[string][]byte) *DBEntity {
	fields := make([]*DBEntityField, 0, len(fTypes))

	for name, value := range fTypes {
		ef := &DBEntityField{
			Name:  name,
			Value: value,
		}
		fields = append(fields, ef)
	}

	return &DBEntity{
		Fields: fields,
	}
}

type DBEntityField struct {
	gorm.Model

	EntityID uint      `gorm:"index"`
	Entity   *DBEntity `gorm:"foreignKey:EntityID;references:ID"`

	Name  string
	Value []byte

	ValueDecoded any `gorm:"-"`
}
