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

func (q *Query) ToDB(db *gorm.DB, fTypes map[string]entities.FieldValueInfo) (*gorm.DB, error) {
	db = db.Select("*")
	db = db.Joins("JOIN db_topic ON db_topic.id = db_entities.id")
	db = db.Where("db_topic.name = ?", q.Topic)

	for _, f := range q.Filters {
		key := f.Field
		keyF := "f_" + f.Field
		keyV := "v_" + f.Field

		db = db.Joins("JOIN db_entity_fields "+keyF+" ON "+keyF+".entity_id = db_entities.id AND "+keyF+".name = ?", key)
		db = db.Joins("JOIN db_entity_field_values " + keyV + " ON " + keyV + ".db_entity_field_id = " + keyF + ".id")
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
		}
	}

	return db, nil
}
