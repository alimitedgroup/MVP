package application

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type ApplyCatalogUpdateService struct {
	applyCatalogUpdatePort port.ApplyCatalogUpdatePort
}

func NewApplyCatalogUpdateService(applyGoodUpdatePort port.ApplyCatalogUpdatePort) *ApplyCatalogUpdateService {
	return &ApplyCatalogUpdateService{applyGoodUpdatePort}
}

func (s *ApplyCatalogUpdateService) ApplyCatalogUpdate(ctx context.Context, cmd port.CatalogUpdateCmd) error {
	good := model.GoodInfo{
		ID:          cmd.GoodId,
		Name:        cmd.Name,
		Description: cmd.Description,
	}

	err := s.applyCatalogUpdatePort.ApplyCatalogUpdate(good)
	if err != nil {
		return err
	}

	return nil
}
