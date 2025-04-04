package persistence

import "errors"

type IStockRepository interface {
	GetStock(goodId string) int64
	SetStock(goodId string, stock int64) bool
	AddStock(goodId string, stock int64) bool
	GetFreeStock(goodId string) int64
	ReserveStock(reservationId string, goodId string, stock int64) error
	UnReserveStock(goodId string, stock int64) error
	GetReservation(reservationId string) (Reservation, error)
}

var ErrNotEnoughGoods = errors.New("not enough goods")
var ErrReservationNotFound = errors.New("reservation not found")

type Reservation struct {
	Goods map[string]int64
}
