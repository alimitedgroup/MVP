package persistence

import "sync"

type StockRepositoryIml struct {
	m           sync.Mutex
	goodToStock map[string]int64
}

func NewStockRepositoryIml() *StockRepositoryIml {
	return &StockRepositoryIml{goodToStock: make(map[string]int64)}
}

func (s *StockRepositoryIml) GetStock(string string) int64 {
	s.m.Lock()
	defer s.m.Unlock()

	return s.goodToStock[string]
}

func (s *StockRepositoryIml) SetStock(string string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, exist := s.goodToStock[string]

	s.goodToStock[string] = stock

	return exist
}

func (s *StockRepositoryIml) AddStock(string string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	prev, exist := s.goodToStock[string]

	s.goodToStock[string] = prev + stock

	return exist
}
