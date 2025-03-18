package port

import (
	"context"
)

type ICreateTransferUseCase interface {
	CreateTransfer(context.Context, CreateTransferCmd) (CreateTransferResponse, error)
}

type CreateTransferCmd struct {
	SenderId   string
	ReceiverId string
	Goods      []CreateTransferGood
}

type CreateTransferGood struct {
	GoodID   string
	Quantity int64
}

type CreateTransferResponse struct {
	TransferID string
}
