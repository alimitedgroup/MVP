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
	UpdateTime    int64
	SenderID      string
	ReceiverID    string
	Goods         []SendTransferUpdateGood
	ReservationID string
}

type SendTransferUpdateGood struct {
	GoodID   string
	Quantity int64
}
