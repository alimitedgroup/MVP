package persistance

import "sync"

type GoodID string

type StockRepository struct {
	sync.Mutex
	goodToStock map[GoodID]int64
}

func NewStockRepository() *StockRepository {
	return &StockRepository{goodToStock: make(map[GoodID]int64)}
}

func (s *StockRepository) GetStock(goodID GoodID) int64 {
	s.Lock()
	defer s.Unlock()

	return s.goodToStock[goodID]
}

func (s *StockRepository) SetStock(goodID GoodID, stock int64) bool {
	s.Lock()
	defer s.Unlock()

	_, exist := s.goodToStock[goodID]

	s.goodToStock[goodID] = stock

	return exist
}

func (s *StockRepository) AddStock(goodID GoodID, stock int64) bool {
	s.Lock()
	defer s.Unlock()

	prev, exist := s.goodToStock[goodID]

	s.goodToStock[goodID] = prev + stock

	return exist
}
