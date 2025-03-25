package port

import "github.com/alimitedgroup/MVP/srv/order/business/model"

type IGetTransferPort interface {
	GetTransfer(model.TransferID) (model.Transfer, error)
	GetAllTransfer() []model.Transfer
}
