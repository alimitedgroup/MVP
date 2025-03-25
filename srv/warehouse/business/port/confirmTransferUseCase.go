package port

import "context"

type IConfirmTransferUseCase interface {
	ConfirmTransfer(context.Context, ConfirmTransferCmd) error
}

type ConfirmTransferCmd struct {
	TransferID    string
	SenderID      string
	ReceiverID    string
	Status        string
	Goods         []TransferUpdateGood
	ReservationID string
}

type TransferUpdateGood struct {
	GoodID   string
	Quantity int64
}
