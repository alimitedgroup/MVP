package persistence

import "github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"

type IGoodRepository interface {
	GetGoods() map[string]catalogCommon.Good
	GetGoodsGlobalQuantity() map[string]int64
	SetGoodQuantity(warehouseID string, goodID string, newQuantity int64) error
	AddGood(goodID string, name string, description string) error
	GetWarehouses() map[string]catalogCommon.Warehouse
}
