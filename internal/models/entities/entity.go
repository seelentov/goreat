package entities

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

var FormatErr = errors.New("format error")

type FieldValuePair struct {
	FieldType FieldType
	Value     any
}

type Entity struct {
	gorm.Model

	TopicID uint   `gorm:"index"`
	Topic   *Topic `gorm:"foreignKey:TopicID;references:ID"`

	Fields []*EntityField `gorm:"foreignKey:EntityID;constraint:OnDelete:CASCADE"`
}

func NewEntity(fTypes map[string]FieldValuePair) *Entity {
	fields := make([]*EntityField, 0, len(fTypes))

	for name, pair := range fTypes {
		ef := NewEntityField(name, pair.FieldType, pair.Value)
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
	*Field

	EntityID uint    `gorm:"index"`
	Entity   *Entity `gorm:"foreignKey:EntityID;references:ID"`

	Value []byte
}

func NewEntityField(name string, fType FieldType, value any) *EntityField {
	v, err := serializeValue(value, fType)
	if err != nil {
		return nil
	}

	return &EntityField{
		Field: NewField(name, fType),
		Value: v,
	}
}

func (ef *EntityField) GetValue() (any, error) {
	return deserializeValue(ef.Value, ef.Field.Type)
}

func (ef *EntityField) SetValue(value any) error {
	v, err := serializeValue(value, ef.Field.Type)
	if err != nil {
		return err
	}
	ef.Value = v
	return nil
}

func serializeValue(data any, fType FieldType) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	switch fType {
	case FieldTypeString:
		str, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("%w: expected string for FieldTypeString, got %T", FormatErr, data)
		}
		return []byte(str), nil

	case FieldTypeInt:
		var val int64
		switch v := data.(type) {
		case int:
			val = int64(v)
		case int32:
			val = int64(v)
		case int64:
			val = v
		case uint:
			val = int64(v)
		case uint32:
			val = int64(v)
		case uint64:
			if v > 1<<63-1 {
				return nil, fmt.Errorf("%w: uint64 value %d overflows int64", FormatErr, v)
			}
			val = int64(v)
		default:
			return nil, fmt.Errorf("%w: expected integer type for FieldTypeInt, got %T", FormatErr, data)
		}
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(val))
		return buf, nil

	case FieldTypeFloat:
		var val float64
		switch v := data.(type) {
		case float32:
			val = float64(v)
		case float64:
			val = v
		default:
			return nil, fmt.Errorf("%w: expected float type for FieldTypeFloat, got %T", FormatErr, data)
		}
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(val))
		return buf, nil

	case FieldTypeBool:
		b, ok := data.(bool)
		if !ok {
			return nil, fmt.Errorf("%w: expected bool for FieldTypeBool, got %T", FormatErr, data)
		}
		if b {
			return []byte{1}, nil
		}
		return []byte{0}, nil

	case FieldTypeDateTime:
		var t time.Time
		switch v := data.(type) {
		case time.Time:
			t = v
		case int64:
			t = time.Unix(0, v)
		case string:
			var err error
			t, err = time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("%w: failed to parse time string: %w", FormatErr, err)
			}
		default:
			return nil, fmt.Errorf("%w: expected time.Time, int64 or string for FieldTypeDateTime, got %T", FormatErr, data)
		}
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(t.UnixNano()))
		return buf, nil

	default:
		return nil, fmt.Errorf("%w: unsupported FieldType: %d", FormatErr, fType)
	}
}

func deserializeValue(data []byte, fType FieldType) (any, error) {
	if data == nil {
		return nil, nil
	}

	switch fType {
	case FieldTypeString:
		return string(data), nil

	case FieldTypeInt:
		if len(data) < 8 {
			return nil, fmt.Errorf("%w: invalid data length for FieldTypeInt: expected 8 bytes, got %d", FormatErr, len(data))
		}
		val := int64(binary.LittleEndian.Uint64(data[:8]))
		return val, nil

	case FieldTypeFloat:
		if len(data) < 8 {
			return nil, fmt.Errorf("%w: invalid data length for FieldTypeFloat: expected 8 bytes, got %d", FormatErr, len(data))
		}
		bits := binary.LittleEndian.Uint64(data[:8])
		return math.Float64frombits(bits), nil

	case FieldTypeBool:
		if len(data) < 1 {
			return nil, fmt.Errorf("%w: invalid data length for FieldTypeBool: expected 1 byte, got %d", FormatErr, len(data))
		}
		return data[0] != 0, nil

	case FieldTypeDateTime:
		if len(data) < 8 {
			return nil, fmt.Errorf("%w: invalid data length for FieldTypeDateTime: expected 8 bytes, got %d", FormatErr, len(data))
		}
		nanos := int64(binary.LittleEndian.Uint64(data[:8]))
		return time.Unix(0, nanos), nil

	default:
		return nil, fmt.Errorf("%w: unsupported FieldType: %d", FormatErr, fType)
	}
}
