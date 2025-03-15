package port

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type IApplyCatalogUpdatePort interface {
	ApplyCatalogUpdate(good model.GoodInfo) error
}
