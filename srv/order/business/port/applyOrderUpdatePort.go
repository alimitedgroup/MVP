package port

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
)

type IApplyOrderUpdatePort interface {
	ApplyOrderUpdate(ApplyOrderUpdateCmd) error
}

type ApplyOrderUpdateCmd struct {
	Id           string
	Status       string
	Name         string
	Email        string
	Address      string
	Goods        []model.GoodStock
	Reservations []string
	UpdateTime   int64
	CreationTime int64
}
