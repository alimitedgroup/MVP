package port

import (
	"context"
	"errors"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IGetTransferUseCase interface {
	GetTransfer(context.Context, GetTransferCmd) (model.Transfer, error)
	GetAllTransfers(context.Context) []model.Transfer
}

type GetTransferCmd string

var ErrTransferNotFound = errors.New("transfer not found")
