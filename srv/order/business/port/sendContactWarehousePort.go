package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type ISendContactWarehousePort interface {
	SendContactWarehouses(context.Context, SendContactWarehouseCmd) error
}

type SendContactWarehouseCmd struct {
	Order                 model.Order
	TransferId            string
	LastContact           int64
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
}
