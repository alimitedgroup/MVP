package application

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type ApplyCatalogUpdateService struct {
	applyCatalogUpdatePort port.IApplyCatalogUpdatePort
}

func NewApplyCatalogUpdateService(applyGoodUpdatePort port.IApplyCatalogUpdatePort) *ApplyCatalogUpdateService {
	return &ApplyCatalogUpdateService{applyGoodUpdatePort}
}

func (s *ApplyCatalogUpdateService) ApplyCatalogUpdate(ctx context.Context, cmd port.CatalogUpdateCmd) error {
	good := model.GoodInfo{
		ID:          model.GoodId(cmd.GoodId),
		Name:        cmd.Name,
		Description: cmd.Description,
	}

	err := s.applyCatalogUpdatePort.ApplyCatalogUpdate(good)
	if err != nil {
		return err
	}

	return nil
}
