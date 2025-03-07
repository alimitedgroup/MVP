package persistence

import "sync"

type StockRepositoryIml struct {
	m           sync.Mutex
	goodToStock map[string]int64
}

func NewStockRepositoryIml() *StockRepositoryIml {
	return &StockRepositoryIml{goodToStock: make(map[string]int64)}
}

func (s *StockRepositoryIml) GetStock(goodId string) int64 {
	s.m.Lock()
	defer s.m.Unlock()

	stock, exist := s.goodToStock[goodId]
	if !exist {
		return 0
	}

	return stock
}

func (s *StockRepositoryIml) SetStock(goodId string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, exist := s.goodToStock[goodId]

	s.goodToStock[goodId] = stock

	return exist
}

func (s *StockRepositoryIml) AddStock(goodId string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	prev, exist := s.goodToStock[goodId]

	s.goodToStock[goodId] = prev + stock

	return exist
}
