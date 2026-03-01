package queries

type OrderDirection int

const (
	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

type Order struct {
	Field     string         `json:"field" binding:"required"`
	Direction OrderDirection `json:"direction" binding:"oneof=0 1"`
}
