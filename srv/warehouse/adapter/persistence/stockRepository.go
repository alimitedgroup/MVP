package persistence

import "errors"

type IStockRepository interface {
	GetStock(goodId string) int64
	SetStock(goodId string, stock int64) bool
	AddStock(goodId string, stock int64) bool
	GetFreeStock(goodId string) int64
	ReserveStock(goodId string, stock int64) error
	UnReserveStock(goodId string, stock int64) error
}

var ErrNotEnoughGoods = errors.New("not enough goods")
