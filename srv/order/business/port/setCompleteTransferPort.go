package port

import "github.com/alimitedgroup/MVP/srv/order/business/model"

type ISetCompleteTransferPort interface {
	SetComplete(model.TransferID) error
	IncrementLinkedStockUpdate(model.TransferID) error
}
