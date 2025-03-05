package port

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type ApplyStockUpdatePort interface {
	ApplyStockUpdate([]model.GoodStock) error
}
