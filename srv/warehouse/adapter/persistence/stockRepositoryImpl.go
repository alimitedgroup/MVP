package persistence

import (
	"sync"
)

type StockRepositoryImpl struct {
	m             sync.Mutex
	goodToStock   map[string]int64
	reservedStock map[string]int64
	reservations  map[string]Reservation
}

func NewStockRepositoryImpl() *StockRepositoryImpl {
	return &StockRepositoryImpl{goodToStock: make(map[string]int64), reservedStock: make(map[string]int64), reservations: make(map[string]Reservation)}
}

func (s *StockRepositoryImpl) GetStock(goodId string) int64 {
	s.m.Lock()
	defer s.m.Unlock()

	stock, exist := s.goodToStock[goodId]
	if !exist {
		return 0
	}

	return stock
}

func (s *StockRepositoryImpl) SetStock(goodId string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, exist := s.goodToStock[goodId]

	s.goodToStock[goodId] = stock

	return exist
}

func (s *StockRepositoryImpl) AddStock(goodId string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	prev, exist := s.goodToStock[goodId]

	s.goodToStock[goodId] = prev + stock

	return exist
}

func (s *StockRepositoryImpl) ReserveStock(reservationId string, goodId string, stock int64) error {
	s.m.Lock()
	defer s.m.Unlock()

	prevReserved, exist := s.reservedStock[goodId]
	if !exist {
		prevReserved = 0
	}

	currStock, exist := s.goodToStock[goodId]
	if !exist {
		currStock = 0
	}

	if currStock-prevReserved < stock {
		return ErrNotEnoughGoods
	}

	s.reservedStock[goodId] = prevReserved + stock

	if _, exist = s.reservations[reservationId]; !exist {
		s.reservations[reservationId] = Reservation{Goods: make(map[string]int64)}
	}

	old, exist := s.reservations[reservationId].Goods[goodId]
	if !exist {
		old = 0
	}
	s.reservations[reservationId].Goods[goodId] = old + stock

	return nil
}

func (s *StockRepositoryImpl) UnReserveStock(goodId string, stock int64) error {
	s.m.Lock()
	defer s.m.Unlock()

	prev, exist := s.reservedStock[goodId]
	if !exist {
		return ErrNotEnoughGoods
	}

	if prev < stock {
		return ErrNotEnoughGoods
	}

	s.reservedStock[goodId] = prev - stock

	return nil
}

func (s *StockRepositoryImpl) GetFreeStock(goodId string) int64 {
	s.m.Lock()
	defer s.m.Unlock()

	stock, exist := s.goodToStock[goodId]
	if !exist {
		return 0
	}

	reserved, exist := s.reservedStock[goodId]
	if !exist {
		return stock
	}

	return stock - reserved
}

func (s *StockRepositoryImpl) GetReservation(reservationId string) (Reservation, error) {
	s.m.Lock()
	defer s.m.Unlock()

	reserv, exist := s.reservations[reservationId]
	if !exist {
		return Reservation{}, ErrReservationNotFound
	}

	return reserv, nil
}
