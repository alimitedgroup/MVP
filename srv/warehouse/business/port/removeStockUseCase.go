package port

import "context"

type IRemoveStockUseCase interface {
	RemoveStock(context.Context, RemoveStockCmd) error
}

type RemoveStockCmd struct {
	GoodID   string
	Quantity int64
}
