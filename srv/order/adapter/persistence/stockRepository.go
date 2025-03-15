package persistence

import "errors"

type IStockRepository interface {
	GetStock(warehouseId string, goodId string) (int64, error)
	SetStock(warehouseId string, goodId string, stock int64) bool
	AddStock(warehouseId string, goodId string, stock int64) (bool, error)
}

var ErrWarehouseNotFound = errors.New("warehouse not found")
var ErrGoodNotFound = errors.New("good not found")
