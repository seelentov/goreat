package queries

import (
	"goreat/internal/models/entities"

	"gorm.io/gorm"
)

type QueryType int

const (
	QueryTypeData QueryType = iota
	QueryTypeCount
	QueryTypeExists
)

type Query struct {
	Topic string

	Filters []*Filter
	Orders  []*Order

	Limit  *uint
	Offset *uint

	Type QueryType
}

func (q *Query) ToDB(db *gorm.DB, fTypes map[string]entities.FieldInfo) (*gorm.DB, error) {
	db = db.Session(&gorm.Session{}).Model(&entities.DBEntity{})
	db = db.Preload("Fields.Value")

	db = db.Joins("JOIN db_topics ON db_topics.id = db_entities.topic_id")
	db = db.Where("db_topics.name = ?", q.Topic)

	joinedFields := make(map[string]bool)

	fields := make([]string, 0, len(q.Filters)+len(q.Orders))
	for _, f := range q.Filters {
		fields = append(fields, f.Field)
	}
	for _, o := range q.Orders {
		fields = append(fields, o.Field)
	}

	for _, key := range fields {
		if !joinedFields[key] {
			keyF := "f_" + key
			keyV := "v_" + key

			db = db.Joins("JOIN db_entity_fields "+keyF+" ON "+keyF+".entity_id = db_entities.id AND "+keyF+".name = ?", key)
			db = db.Joins("JOIN db_entity_field_values " + keyV + " ON " + keyV + ".db_entity_field_id = " + keyF + ".id")
			joinedFields[key] = true
		}
	}

	for _, f := range q.Filters {
		fType := fTypes[f.Field]
		keyV := "v_" + f.Field

		var column string
		switch fType.FieldType {
		case entities.FieldTypeInt:
			column = keyV + ".value_int"
		case entities.FieldTypeString:
			column = keyV + ".value_str"
		case entities.FieldTypeFloat:
			column = keyV + ".value_float"
		case entities.FieldTypeBool:
			column = keyV + ".value_bool"
		case entities.FieldTypeDateTime:
			column = keyV + ".value_time"
		default:
			return nil, entities.ErrType
		}

		switch f.Type {
		case FilterTypeEquals:
			db = db.Where(column+" = ?", f.Value)
		case FilterTypeNotEquals:
			db = db.Where(column+" != ?", f.Value)
		case FilterTypeContains:
			db = db.Where(column+" LIKE ?", "%"+f.Value+"%")
		case FilterTypeGreaterThan:
			db = db.Where(column+" > ?", f.Value)
		case FilterTypeLessThan:
			db = db.Where(column+" < ?", f.Value)
		case FilterTypeGreaterThanOrEquals:
			db = db.Where(column+" >= ?", f.Value)
		case FilterTypeLessThanOrEquals:
			db = db.Where(column+" <= ?", f.Value)
		default:
			return nil, entities.ErrType
		}
	}

	for _, o := range q.Orders {
		fType := fTypes[o.Field]
		keyV := "v_" + o.Field

		var column string
		switch fType.FieldType {
		case entities.FieldTypeInt:
			column = keyV + ".value_int"
		case entities.FieldTypeString:
			column = keyV + ".value_str"
		case entities.FieldTypeFloat:
			column = keyV + ".value_float"
		case entities.FieldTypeBool:
			column = keyV + ".value_bool"
		case entities.FieldTypeDateTime:
			column = keyV + ".value_time"
		default:
			return nil, entities.ErrType
		}

		dir := "ASC"
		if o.Direction == OrderDirectionDesc {
			dir = "DESC"
		}

		db = db.Order(column + " " + dir)
	}

	if q.Limit != nil {
		db = db.Limit(int(*q.Limit))
	}

	if q.Offset != nil {
		db = db.Offset(int(*q.Offset))
	}

	if q.Type == QueryTypeExists {
		db = db.Limit(1)
	}

	return db, nil
}
