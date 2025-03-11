package port

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type IApplyStockUpdatePort interface {
	ApplyStockUpdate([]model.GoodStock) error
}
