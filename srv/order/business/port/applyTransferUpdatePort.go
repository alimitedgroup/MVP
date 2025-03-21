package port

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IApplyTransferUpdatePort interface {
	ApplyTransferUpdate(ApplyTransferUpdateCmd)
}

type ApplyTransferUpdateCmd struct {
	Id            string
	Status        string
	SenderId      string
	ReceiverId    string
	Goods         []model.GoodStock
	ReservationId string
	UpdateTime    int64
	CreationTime  int64
}
