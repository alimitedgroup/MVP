package port

type IApplyStockUpdateUseCase interface {
	ApplyStockUpdate(StockUpdateCmd)
}

type StockUpdateCmd struct {
	ID            string
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
