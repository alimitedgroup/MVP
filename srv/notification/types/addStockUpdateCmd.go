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

func NewAddStockUpdateCmd(warehouseID, updateType, orderID, transferID string, goods []StockGood, timestamp int64) *AddStockUpdateCmd {
	return &AddStockUpdateCmd{
		WarehouseID: warehouseID,
		Type:        updateType,
		Goods:       goods,
		OrderID:     orderID,
		TransferID:  transferID,
		Timestamp:   timestamp,
	}
}

func (a *AddStockUpdateCmd) GetWarehouseID() string { return a.WarehouseID }
func (a *AddStockUpdateCmd) GetGoods() []StockGood  { return a.Goods }
