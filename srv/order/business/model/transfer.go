package model

type TransferID string

type Transfer struct {
	Id                TransferID
	SenderId          WarehouseID
	ReceiverId        WarehouseID
	Status            string
	UpdateTime        int64
	CreationTime      int64
	LinkedStockUpdate int
	ReservationID     string
	Goods             []GoodStock
}
