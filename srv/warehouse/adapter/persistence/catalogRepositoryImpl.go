package persistence

import "sync"

type CatalogRepositoryImpl struct {
	m       sync.Mutex
	goodMap map[string]Good
}

func NewCatalogRepositoryImpl() *CatalogRepositoryImpl {
	return &CatalogRepositoryImpl{goodMap: make(map[string]Good)}
}

func (s *CatalogRepositoryImpl) GetGood(goodId string) *Good {
	s.m.Lock()
	defer s.m.Unlock()

	good, exist := s.goodMap[goodId]
	if !exist {
		return nil
	}

	return &good
}

func (s *CatalogRepositoryImpl) SetGood(goodId string, name string, description string) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, exist := s.goodMap[goodId]

	s.goodMap[goodId] = Good{
		Id:          goodId,
		Name:        name,
		Description: description,
	}

	return exist
}
