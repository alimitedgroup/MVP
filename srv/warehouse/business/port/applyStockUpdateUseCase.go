package port

type IApplyStockUpdateUseCase interface {
	ApplyStockUpdate(StockUpdateCmd)
}

type StockUpdateCmd struct {
	ID            string
	Type          StockUpdateCmdType
	Goods         []StockUpdateCmdGood
	OrderID       string
	TransferID    string
	ReservationID string
	Timestamp     int64
}

type StockUpdateCmdType string

const (
	StockUpdateCmdTypeAdd      StockUpdateCmdType = "add"
	StockUpdateCmdTypeRemove   StockUpdateCmdType = "remove"
	StockUpdateCmdTypeOrder    StockUpdateCmdType = "order"
	StockUpdateCmdTypeTransfer StockUpdateCmdType = "transfer"
)

type StockUpdateCmdGood struct {
	GoodID   string
	Quantity int64
	Delta    int64
}
