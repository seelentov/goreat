package queries

type QueryType int

const (
	QueryTypeData QueryType = iota
	QueryTypeCount
	QueryTypeExists
)

type Query struct {
	Filters []*Filter
	Orders  []*Order

	Limit  *uint
	Offset *uint

	Type QueryType
}
