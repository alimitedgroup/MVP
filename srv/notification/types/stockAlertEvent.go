package types

type StockAlertEvent struct {
	Id              string
	Status          StockStatus
	GoodID          string
	CurrentQuantity int
	Operator        string
	Threshold       int
	Timestamp       int64
	RuleId          string
}

type StockStatus string

var (
	StockPending      StockStatus = "Pending"
	StockAcknowledged StockStatus = "Acknowledged"
)
