package port

import "context"

type IApplyStockUpdateUseCase interface {
	ApplyStockUpdate(context.Context, StockUpdateCmd) error
}

type StockUpdateCmd struct {
	ID          string
	WarehouseID string
	Type        StockUpdateCmdType
	Goods       []StockUpdateCmdGood
	OrderID     string
	TransferID  string
	Timestamp   int64
}

type StockUpdateCmdType string

const (
	StockUpdateCmdTypeAdd    StockUpdateCmdType = "add"
	StockUpdateCmdTypeRemove StockUpdateCmdType = "remove"
)

type StockUpdateCmdGood struct {
	GoodID   string
	Quantity int64
	Delta    int64
}
