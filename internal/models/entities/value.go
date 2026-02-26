package entities

import (
	"errors"
	"time"
)

var ErrType = errors.New("type not implemented")

type Value struct {
	ValueInt   *int64
	ValueFloat *float64
	ValueStr   *string
	ValueBool  *bool
	ValueTime  *time.Time
}

func NewValue(d any) (*Value, error) {
	val := &Value{}

	switch v := d.(type) {
	case int:
		va := int64(v)
		val.ValueInt = &va
	case int8:
		va := int64(v)
		val.ValueInt = &va
	case int16:
		va := int64(v)
		val.ValueInt = &va
	case int32:
		va := int64(v)
		val.ValueInt = &va
	case int64:
		val.ValueInt = &v

	case float32:
		va := float64(v)
		val.ValueFloat = &va
	case float64:
		val.ValueFloat = &v

	case bool:
		val.ValueBool = &v
	case string:
		val.ValueStr = &v

	case time.Time:
		val.ValueTime = &v

	default:
		return nil, ErrType
	}

	return val, nil
}

func (v *Value) GetValue() any {
	if v.ValueBool != nil {
		return *v.ValueBool
	} else if v.ValueInt != nil {
		return *v.ValueInt
	} else if v.ValueFloat != nil {
		return *v.ValueFloat
	} else if v.ValueStr != nil {
		return *v.ValueStr
	} else if v.ValueTime != nil {
		return *v.ValueTime
	}

	return nil
}

type DBEntityFieldValue struct {
	Value

	DBEntityFieldID uint
	DBEntityField   *DBEntityField `gorm:"foreignKey:DBEntityFieldID"`
}

func NewDBEntityFieldValue(d any) (*DBEntityFieldValue, error) {
	val, err := NewValue(d)
	if err != nil {
		return nil, err
	}

	return &DBEntityFieldValue{
		Value: *val,
	}, nil
}
