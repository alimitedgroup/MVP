package portin

import "github.com/alimitedgroup/MVP/common/dto"

type Order interface {
	CreateOrder(any) (string, error)
	GetOrders() ([]dto.Order, error)
	CreateTransfer(any) (string, error)
	GetTransfers() ([]dto.Transfer, error)
}
