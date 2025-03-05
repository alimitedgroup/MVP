package port

import "context"

type ApplyStockUpdateUseCase interface {
	ApplyStockUpdate(context.Context, StockUpdateCmd) error
}

type StockUpdateCmd struct {
	ID         string
	Type       string
	Goods      []StockUpdateCmdGood
	OrderID    string
	TransferID string
	Timestamp  int64
}

type StockUpdateCmdGood struct {
	GoodID   string
	Quantity int64
	Delta    int64
}
