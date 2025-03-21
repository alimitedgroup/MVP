package port

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type IApplyStockUpdatePort interface {
	ApplyStockUpdate([]model.GoodStock)
}
