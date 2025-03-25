package port

import "context"

type IApplyTransferUpdateUseCase interface {
	ApplyTransferUpdate(context.Context, TransferUpdateCmd)
}

type TransferUpdateCmd struct {
	ID            string
	SenderID      string
	ReceiverID    string
	Goods         []TransferUpdateGood
	ReservationID string
	Status        string
	CreationTime  int64
	UpdateTime    int64
}

type TransferUpdateGood struct {
	GoodID   string
	Quantity int64
}
