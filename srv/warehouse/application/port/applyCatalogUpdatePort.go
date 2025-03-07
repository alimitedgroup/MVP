package port

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type ApplyCatalogUpdatePort interface {
	ApplyCatalogUpdate(good model.GoodInfo) error
}
