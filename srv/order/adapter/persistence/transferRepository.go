package persistence

import (
	"errors"
)

type ITransferRepository interface {
	GetTransfer(transferId string) (Transfer, error)
	GetTransfers() []Transfer
	SetTransfer(transferId string, transfer Transfer) bool
	SetComplete(transferId string) error
	IncrementLinkedStockUpdate(transferId string) error
}

type Transfer struct {
	ID                string
	Status            string
	SenderID          string
	ReceiverID        string
	Goods             []TransferUpdateGood
	LinkedStockUpdate int
	ReservationID     string
	CreationTime      int64
	UpdateTime        int64
}

type TransferUpdateGood struct {
	GoodID   string
	Quantity int64
}

var ErrTransferNotFound = errors.New("transfer not found")
