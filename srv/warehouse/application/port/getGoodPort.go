package port

import "github.com/alimitedgroup/MVP/srv/warehouse/model"

type GetGoodPort interface {
	GetGood(goodId string) *model.GoodInfo
}
