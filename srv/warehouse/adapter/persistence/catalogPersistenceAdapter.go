package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type CatalogPersistanceAdapter struct {
	catalogRepo ICatalogRepository
}

func NewCatalogPersistanceAdapter(catalogRepo ICatalogRepository) *CatalogPersistanceAdapter {
	return &CatalogPersistanceAdapter{catalogRepo}
}

func (s *CatalogPersistanceAdapter) ApplyCatalogUpdate(good model.GoodInfo) error {
	s.catalogRepo.SetGood(string(good.ID), good.Name, good.Description)

	return nil
}

func (s *CatalogPersistanceAdapter) GetGood(goodId model.GoodID) *model.GoodInfo {
	good := s.catalogRepo.GetGood(string(goodId))
	if good == nil {
		return nil
	}

	return &model.GoodInfo{
		ID:          model.GoodID(good.Id),
		Name:        good.Name,
		Description: good.Description,
	}
}
