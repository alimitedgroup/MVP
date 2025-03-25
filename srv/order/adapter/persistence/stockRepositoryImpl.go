package persistence

import (
	"sync"
)

type WarehouseStock struct {
	goodToStock map[string]int64
}

type StockRepositoryImpl struct {
	m              sync.Mutex
	warehouseMap   map[string]WarehouseStock
	globalStockMap map[string]int64
}

func NewStockRepositoryImpl() *StockRepositoryImpl {
	return &StockRepositoryImpl{warehouseMap: make(map[string]WarehouseStock), globalStockMap: make(map[string]int64)}
}

func (s *StockRepositoryImpl) GetStock(warehouseId string, goodId string) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	warehouse, exist := s.warehouseMap[warehouseId]
	if !exist {
		return 0, ErrWarehouseNotFound
	}

	stock, exist := warehouse.goodToStock[goodId]
	if !exist {
		return 0, ErrGoodNotFound
	}

	return stock, nil
}

func (s *StockRepositoryImpl) SetStock(warehouseId string, goodId string, stock int64) bool {
	s.m.Lock()
	defer s.m.Unlock()

	warehouse, exist := s.warehouseMap[warehouseId]
	if !exist {
		warehouse = WarehouseStock{goodToStock: make(map[string]int64)}
		s.warehouseMap[warehouseId] = warehouse
	} else {
		_, exist = warehouse.goodToStock[goodId]
	}

	warehouse.goodToStock[goodId] = stock

	prevGlobalStock, globalExist := s.globalStockMap[goodId]
	if !globalExist {
		prevGlobalStock = 0
	}
	s.globalStockMap[goodId] = prevGlobalStock + stock

	return exist
}

func (s *StockRepositoryImpl) AddStock(warehouseId string, goodId string, stock int64) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	warehouse, exist := s.warehouseMap[warehouseId]
	if !exist {
		return false, ErrWarehouseNotFound
	}

	prevStock, exist := warehouse.goodToStock[goodId]
	if !exist {
		return false, ErrGoodNotFound
	}

	warehouse.goodToStock[goodId] = prevStock + stock

	prevGlobalStock, globalExist := s.globalStockMap[goodId]
	if !globalExist {
		prevGlobalStock = 0
	}
	s.globalStockMap[goodId] = prevGlobalStock + stock

	return exist, nil
}

func (s *StockRepositoryImpl) GetGlobalStock(goodId string) int64 {
	s.m.Lock()
	defer s.m.Unlock()

	stock, exist := s.globalStockMap[goodId]
	if !exist {
		return 0
	}

	return stock
}

func (s *StockRepositoryImpl) GetWarehouses() []string {
	s.m.Lock()
	defer s.m.Unlock()

	warehouses := make([]string, 0, len(s.warehouseMap))
	for warehouseId := range s.warehouseMap {
		warehouses = append(warehouses, warehouseId)
	}

	return warehouses
}
