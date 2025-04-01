package types

type AddStockUpdateCmd struct {
	WarehouseID string
	Type        string
	Goods       []StockGood
	OrderID     string
	TransferID  string
	Timestamp   int64
}

type StockGood struct {
	ID       string
	Quantity int
	Delta    int
}
