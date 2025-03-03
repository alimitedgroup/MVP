package application

import "github.com/alimitedgroup/MVP/srv/warehouse/application/port"

type UpdateStockService struct {
}

func NewUpdateStockService() *UpdateStockService {
	return &UpdateStockService{}
}

func (s *UpdateStockService) UpdateStock(port.UpdateStockCommand) error {
	return nil
}
