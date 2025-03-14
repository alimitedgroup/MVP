package port

import "github.com/alimitedgroup/MVP/srv/warehouse/model"

type IGetGoodPort interface {
	GetGood(goodId model.GoodId) *model.GoodInfo
}
