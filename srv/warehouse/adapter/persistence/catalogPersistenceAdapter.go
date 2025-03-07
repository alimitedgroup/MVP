package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type CatalogPersistanceAdapter struct {
	catalogRepo CatalogRepository
}

func NewCatalogPersistanceAdapter(catalogRepo CatalogRepository) *CatalogPersistanceAdapter {
	return &CatalogPersistanceAdapter{catalogRepo}
}

func (s *CatalogPersistanceAdapter) ApplyCatalogUpdate(good model.GoodInfo) error {
	s.catalogRepo.SetGood(good.ID, good.Name, good.Description)

	return nil
}

func (s *CatalogPersistanceAdapter) GetGood(goodId string) *model.GoodInfo {
	good := s.catalogRepo.GetGood(goodId)
	if good == nil {
		return nil
	}

	return &model.GoodInfo{
		ID:          good.Id,
		Name:        good.Name,
		Description: good.Description,
	}
}
