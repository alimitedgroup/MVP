package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type ISendContactWarehousePort interface {
	SendContactWarehouses(context.Context, SendContactWarehouseCmd) error
}

type SendContactWarehouseCmd struct {
	Order                 *model.Order
	Transfer              *model.Transfer
	Type                  SendContactWarehouseType
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
	RetryInTime           int64
	RetryUntil            int64
}

type SendContactWarehouseType string

var (
	SendContactWarehouseTypeOrder    SendContactWarehouseType = "order"
	SendContactWarehouseTypeTransfer SendContactWarehouseType = "transfer"
)
