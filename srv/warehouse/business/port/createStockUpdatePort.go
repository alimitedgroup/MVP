package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type ICreateStockUpdatePort interface {
	CreateStockUpdate(context.Context, CreateStockUpdateCmd) error
}

type CreateStockUpdateCmd struct {
	Type          CreateStockUpdateCmdType
	Goods         []CreateStockUpdateCmdGood
	OrderID       string
	TransferID    string
	ReservationID string
}

type CreateStockUpdateCmdType string

const (
	CreateStockUpdateCmdTypeAdd      CreateStockUpdateCmdType = "add"
	CreateStockUpdateCmdTypeRemove   CreateStockUpdateCmdType = "remove"
	CreateStockUpdateCmdTypeOrder    CreateStockUpdateCmdType = "order"
	CreateStockUpdateCmdTypeTransfer CreateStockUpdateCmdType = "transfer"
)

type CreateStockUpdateCmdGood struct {
	Good         model.GoodStock
	QuantityDiff int64
}
