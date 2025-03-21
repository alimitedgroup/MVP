package model

type TransferID string

type Transfer struct {
	ID                string
	SenderId          string
	ReceiverId        string
	Status            string
	UpdateTime        int64
	CreationTime      int64
	LinkedStockUpdate int
	ReservationID     string
	Goods             []GoodStock
}
