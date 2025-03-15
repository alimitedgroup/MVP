package port

import "context"

type ISendContactWarehousePort interface {
	SendContactWarehouses(context.Context, SendContactWarehouseCmd) error
}

type SendContactWarehouseCmd struct {
	OrderId               string
	TransferId            string
	LastContact           int64
	ConfirmedReservations []ConfirmedReservation
	ExcludeWarehouses     []string
}
