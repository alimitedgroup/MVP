package port

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type CreateStockUpdatePort interface {
	CreateStockUpdate(context.Context, CreateStockUpdateCmd) error
}

type CreateStockUpdateCmd struct {
	Type  CreateStockUpdateCmdType
	Goods []CreateStockUpdateCmdGood
}

type CreateStockUpdateCmdType string

const (
	CreateStockUpdateCmdTypeAdd    CreateStockUpdateCmdType = "add"
	CreateStockUpdateCmdTypeRemove CreateStockUpdateCmdType = "remove"
)

type CreateStockUpdateCmdGood struct {
	Good         model.GoodStock
	QuantityDiff int64
}
