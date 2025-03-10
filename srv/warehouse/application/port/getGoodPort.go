package port

import "github.com/alimitedgroup/MVP/srv/warehouse/model"

type IGetGoodPort interface {
	GetGood(goodId string) *model.GoodInfo
}
