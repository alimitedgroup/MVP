package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type ICreateStockUpdatePort interface {
	CreateStockUpdate(context.Context, CreateStockUpdateCmd) error
}

type CreateStockUpdateCmd struct {
	Type          CreateStockUpdateType
	Goods         []CreateStockUpdateGood
	OrderID       string
	TransferID    string
	ReservationID string
}

type CreateStockUpdateType string

const (
	CreateStockUpdateCmdTypeAdd      CreateStockUpdateType = "add"
	CreateStockUpdateCmdTypeRemove   CreateStockUpdateType = "remove"
	CreateStockUpdateCmdTypeOrder    CreateStockUpdateType = "order"
	CreateStockUpdateCmdTypeTransfer CreateStockUpdateType = "transfer"
)

type CreateStockUpdateGood struct {
	Good         model.GoodStock
	QuantityDiff int64
}
