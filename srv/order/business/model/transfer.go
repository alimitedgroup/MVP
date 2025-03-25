package model

type TransferID string

type Transfer struct {
	ID                string
	SenderID          string
	ReceiverID        string
	Status            string
	UpdateTime        int64
	CreationTime      int64
	LinkedStockUpdate int
	ReservationID     string
	Goods             []GoodStock
}
