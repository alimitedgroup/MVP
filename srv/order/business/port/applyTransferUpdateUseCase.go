package port

import "context"

type IApplyTransferUpdateUseCase interface {
	ApplyTransferUpdate(context.Context, TransferUpdateCmd)
}

type TransferUpdateCmd struct {
	ID            string
	SenderId      string
	ReceiverId    string
	Goods         []TransferUpdateGood
	ReservationId string
	Status        string
	CreationTime  int64
	UpdateTime    int64
}

type TransferUpdateGood struct {
	GoodID   string
	Quantity int64
}
