package portout

import (
	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
)

type OrderPortOut interface {
	CreateOrder(dto request.CreateOrderRequestDTO) (response.OrderCreateInfo, error)
	GetAllOrders() ([]response.OrderInfo, error)
	CreateTransfer(dto request.CreateTransferRequestDTO) (response.TransferCreateInfo, error)
	GetAllTransfers() ([]response.TransferInfo, error)
}
