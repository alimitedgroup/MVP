package application

import (
	"log"

	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
)

type UpdateStockService struct {
	saveStockUpdatePort port.SaveUpdateStockPort
}

func NewUpdateStockService(saveStockUpdatePort port.SaveUpdateStockPort) *UpdateStockService {
	return &UpdateStockService{saveStockUpdatePort}
}

func (s *UpdateStockService) UpdateStock(cmd port.UpdateStockCmd) error {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       good.GoodID,
			Quantity: good.Quantity,
		})
	}

	err := s.saveStockUpdatePort.SaveUpdateStock(goods)
	if err != nil {
		return err
	}

	log.Printf("loaded stock update %s\n", cmd.ID)

	return nil
}
