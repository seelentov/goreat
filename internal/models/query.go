package models

import "gorm.io/gorm"

type OrderDirection int

const (
	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

type Order struct {
	Field     string
	Direction OrderDirection
}

type FilterType int

const (
	FilterTypeEquals FilterType = iota
	FilterTypeNotEquals
	FilterTypeContains
	FilterTypeGreaterThan
	FilterTypeLessThan
	FilterTypeGreaterThanOrEquals
	FilterTypeLessThanOrEquals
)

type Filters struct {
	Field string
	Type  FilterType
	Value string
}

type QueryType int

const (
	QueryTypeData QueryType = iota
	QueryTypeCount
	QueryTypeExists
)

type QueryData struct {
	gorm.Model

	PageID uint
	Page   *Page
}

type Query struct {
	Filters []*Filters
	Orders  []*Order

	Limit  uint
	Offset uint
}
