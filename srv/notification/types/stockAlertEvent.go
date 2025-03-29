package types

type StockAlertEvent struct {
	GoodID          string
	CurrentQuantity int
	Operator        string
	Threshold       int
	Timestamp       int64
}
