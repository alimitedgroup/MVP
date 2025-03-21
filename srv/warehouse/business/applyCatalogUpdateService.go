package business

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
)

type ApplyCatalogUpdateService struct {
	applyCatalogUpdatePort port.IApplyCatalogUpdatePort
}

func NewApplyCatalogUpdateService(applyGoodUpdatePort port.IApplyCatalogUpdatePort) *ApplyCatalogUpdateService {
	return &ApplyCatalogUpdateService{applyGoodUpdatePort}
}

func (s *ApplyCatalogUpdateService) ApplyCatalogUpdate(cmd port.CatalogUpdateCmd) {
	good := model.GoodInfo{
		ID:          cmd.GoodID,
		Name:        cmd.Name,
		Description: cmd.Description,
	}
	s.applyCatalogUpdatePort.ApplyCatalogUpdate(good)
}
