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
	Field string
	Type  FilterType
	Value string
}
