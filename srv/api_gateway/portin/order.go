package portin

import "github.com/alimitedgroup/MVP/common/dto"

type Order interface {
	CreateOrder(name string, fullname string, address string, goods map[string]int64) (string, error)
	GetOrders() ([]dto.Order, error)
	CreateTransfer(senderID string, receiverID string, goods map[string]int64) (string, error)
	GetTransfers() ([]dto.Transfer, error)
}
