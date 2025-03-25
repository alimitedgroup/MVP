package port

import "github.com/alimitedgroup/MVP/srv/warehouse/business/model"

type IGetGoodPort interface {
	GetGood(goodId model.GoodID) *model.GoodInfo
}
