package port

type UpdateStockUseCase interface {
	UpdateStock(UpdateStockCmd) error
}

type UpdateStockCmd struct {
	ID         string
	Type       string
	Goods      []UpdateStockCommandGood
	OrderID    string
	TransferID string
	Timestamp  int64
}

type UpdateStockCommandGood struct {
	GoodID   string
	Quantity int64
	Delta    int64
}
