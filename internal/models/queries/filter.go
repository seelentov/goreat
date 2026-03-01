package queries

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

type Filter struct {
	Field string     `json:"field,omitempty" binding:"required"`
	Type  FilterType `json:"type" binding:"oneof=0 1 2 3 4 5 6"`
	Value string     `json:"value,omitempty" binding:"required"`
}
