package persistence

import (
	"github.com/alimitedgroup/MVP/common/dto"
)

type IGoodRepository interface {
	GetGoodsGlobalQuantity() map[string]int64
	SetGoodQuantity(warehouseID string, goodID string, newQuantity int64) error
	GetWarehouses() map[string]dto.Warehouse
}
