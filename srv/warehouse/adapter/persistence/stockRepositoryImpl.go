package persistence

import "sync"

type StockRepositoryImpl struct {
	m           sync.Mutex
	goodToStock map[string]int64
}

func NewStockRepositoryImpl() *StockRepositoryImpl {
	return &StockRepositoryImpl{goodToStock: make(map[string]int64)}
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
