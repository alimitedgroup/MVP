package port

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IApplyOrderUpdatePort interface {
	ApplyOrderUpdate(ApplyOrderUpdateCmd)
}

type ApplyOrderUpdateCmd struct {
	ID           string
	Status       string
	Name         string
	FullName     string
	Address      string
	Goods        []model.GoodStock
	Reservations []string
	UpdateTime   int64
	CreationTime int64
}
