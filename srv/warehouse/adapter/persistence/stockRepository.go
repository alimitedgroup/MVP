package persistence

import "sync"

type StockRepository struct {
	m           sync.Mutex
	goodToStock map[string]int64
}

func NewStockRepository() *StockRepository {
	return &StockRepository{goodToStock: make(map[string]int64)}
}

func (s *StockRepository) GetStock(string string) int64 {
	s.m.Lock()
	defer s.m.Unlock()

	return s.goodToStock[string]
}

func (s *StockRepository) SetStock(string string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, exist := s.goodToStock[string]

	s.goodToStock[string] = stock

	return exist
}

func (s *StockRepository) AddStock(string string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	prev, exist := s.goodToStock[string]

	s.goodToStock[string] = prev + stock

	return exist
}
