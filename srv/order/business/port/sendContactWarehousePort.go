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
	LastContact           int64
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
}

type SendContactWarehouseType string

var (
	SendContactWarehouseTypeOrder    SendContactWarehouseType = "order"
	SendContactWarehouseTypeTransfer SendContactWarehouseType = "transfer"
)
