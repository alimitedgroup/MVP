package persistence

import (
	"github.com/alimitedgroup/MVP/common/dto"
)

type IGoodRepository interface {
	GetGoods() map[string]dto.Good
	GetGoodsGlobalQuantity() map[string]int64
	SetGoodQuantity(warehouseID string, goodID string, newQuantity int64) error
	AddGood(goodID string, name string, description string) error
	GetWarehouses() map[string]dto.Warehouse
}
