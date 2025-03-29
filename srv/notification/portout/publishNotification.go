package portout

type IStockEventPublisher interface {
	PublishStockAlert(alert StockAlertEvent) error
}

type StockAlertEvent struct {
	GoodID          string
	CurrentQuantity int
	Operator        string
	Threshold       int
	Timestamp       int64
}
