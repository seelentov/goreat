package entities

import (
	"reflect"

	"gorm.io/gorm"
)

type DBEntityField struct {
	gorm.Model

	Field

	EntityID uint
	Entity   *DBEntity `gorm:"foreignKey:EntityID"`

	Value []*DBEntityFieldValue `gorm:"foreignKey:DBEntityFieldID;constraint:OnDelete:CASCADE"`
}

func NewDBEntityField(name string, value any) (*DBEntityField, error) {
	d := &DBEntityField{}
	d.Name = name

	if err := d.SetValue(value); err != nil {
		return nil, err
	}

	return d, nil
}

func (e *DBEntityField) SetValue(value any) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Slice {
		size := rv.Len()
		e.Value = make([]*DBEntityFieldValue, size)
		for i := 0; i < size; i++ {
			val, err := NewDBEntityFieldValue(rv.Index(i).Interface())
			if err != nil {
				return err
			}
			e.Value[i] = val
		}
	} else {
		val, err := NewDBEntityFieldValue(value)
		if err != nil {
			return err
		}
		e.Value = []*DBEntityFieldValue{val}
	}

	return nil
}
