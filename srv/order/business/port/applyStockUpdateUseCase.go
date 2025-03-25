package port

import "context"

type IApplyStockUpdateUseCase interface {
	ApplyStockUpdate(context.Context, StockUpdateCmd) error
}

type StockUpdateCmd struct {
	ID            string
	WarehouseID   string
	Type          StockUpdateType
	Goods         []StockUpdateGood
	OrderID       string
	TransferID    string
	ReservationID string
	Timestamp     int64
}

type StockUpdateType string

const (
	StockUpdateCmdTypeAdd      StockUpdateType = "add"
	StockUpdateCmdTypeRemove   StockUpdateType = "remove"
	StockUpdateCmdTypeOrder    StockUpdateType = "order"
	StockUpdateCmdTypeTransfer StockUpdateType = "transfer"
)

type StockUpdateGood struct {
	GoodID   string
	Quantity int64
	Delta    int64
}
