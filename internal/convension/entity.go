package convension

import (
	"encoding/binary"
	"errors"
	"fmt"
	"goreat/internal/models/entities"
	"math"
	"time"
)

var ErrFormat = errors.New("format error")

func SerializeValue(data any, fType entities.FieldType, cType entities.ContainerType) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	switch fType {
	case entities.FieldTypeString:
		str, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("%w: expected string for FieldTypeString, got %T", ErrFormat, data)
		}
		return []byte(str), nil

	case entities.FieldTypeInt:
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
				return nil, fmt.Errorf("%w: uint64 value %d overflows int64", ErrFormat, v)
			}
			val = int64(v)
		default:
			return nil, fmt.Errorf("%w: expected integer type for FieldTypeInt, got %T", ErrFormat, data)
		}
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(val))
		return buf, nil

	case entities.FieldTypeFloat:
		var val float64
		switch v := data.(type) {
		case float32:
			val = float64(v)
		case float64:
			val = v
		default:
			return nil, fmt.Errorf("%w: expected float type for FieldTypeFloat, got %T", ErrFormat, data)
		}
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, math.Float64bits(val))
		return buf, nil

	case entities.FieldTypeBool:
		b, ok := data.(bool)
		if !ok {
			return nil, fmt.Errorf("%w: expected bool for FieldTypeBool, got %T", ErrFormat, data)
		}
		if b {
			return []byte{1}, nil
		}
		return []byte{0}, nil

	case entities.FieldTypeDateTime:
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
				return nil, fmt.Errorf("%w: failed to parse time string: %w", ErrFormat, err)
			}
		default:
			return nil, fmt.Errorf("%w: expected time.Time, int64 or string for FieldTypeDateTime, got %T", ErrFormat, data)
		}
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(t.UnixNano()))
		return buf, nil

	default:
		return nil, fmt.Errorf("%w: unsupported FieldType: %d", ErrFormat, fType)
	}
}

func DeserializeValue(data []byte, fType entities.FieldType, cType entities.ContainerType) (any, error) {
	if data == nil {
		return nil, nil
	}

	switch fType {
	case entities.FieldTypeString:
		return string(data), nil

	case entities.FieldTypeInt:
		if len(data) < 8 {
			return nil, fmt.Errorf("%w: invalid data length for FieldTypeInt: expected 8 bytes, got %d", ErrFormat, len(data))
		}
		val := int64(binary.LittleEndian.Uint64(data[:8]))
		return val, nil

	case entities.FieldTypeFloat:
		if len(data) < 8 {
			return nil, fmt.Errorf("%w: invalid data length for FieldTypeFloat: expected 8 bytes, got %d", ErrFormat, len(data))
		}
		bits := binary.LittleEndian.Uint64(data[:8])
		return math.Float64frombits(bits), nil

	case entities.FieldTypeBool:
		if len(data) < 1 {
			return nil, fmt.Errorf("%w: invalid data length for FieldTypeBool: expected 1 byte, got %d", ErrFormat, len(data))
		}
		return data[0] != 0, nil

	case entities.FieldTypeDateTime:
		if len(data) < 8 {
			return nil, fmt.Errorf("%w: invalid data length for FieldTypeDateTime: expected 8 bytes, got %d", ErrFormat, len(data))
		}
		nanos := int64(binary.LittleEndian.Uint64(data[:8]))
		return time.Unix(0, nanos), nil

	default:
		return nil, fmt.Errorf("%w: unsupported FieldType: %d", ErrFormat, fType)
	}
}
