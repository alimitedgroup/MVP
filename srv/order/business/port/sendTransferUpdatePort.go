package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type ISendTransferUpdatePort interface {
	SendTransferUpdate(context.Context, SendTransferUpdateCmd) (model.Transfer, error)
}

type SendTransferUpdateCmd struct {
	ID            string
	Status        string
	CreationTime  int64
	SenderId      string
	ReceiverId    string
	Goods         []SendTransferUpdateGood
	ReservationId string
}

type SendTransferUpdateGood struct {
	GoodID   string
	Quantity int64
}
