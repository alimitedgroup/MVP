package types

type StockAlertEvent struct {
	Id              string
	Status          string
	GoodID          string
	CurrentQuantity int
	Operator        string
	Threshold       int
	Timestamp       int64
}
