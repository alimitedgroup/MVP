package port

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IApplyTransferUpdatePort interface {
	ApplyTransferUpdate(ApplyTransferUpdateCmd)
}

type ApplyTransferUpdateCmd struct {
	ID            string
	Status        string
	SenderID      string
	ReceiverID    string
	Goods         []model.GoodStock
	ReservationID string
	UpdateTime    int64
	CreationTime  int64
}
