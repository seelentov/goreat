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
	panic("implement me")
}
