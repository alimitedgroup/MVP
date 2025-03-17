package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
)

type ApplyStockUpdateService struct {
	applyStockUpdatePort port.IApplyStockUpdatePort
	idempotentPort       port.IIdempotentPort
}

func NewApplyStockUpdateService(applyStockUpdatePort port.IApplyStockUpdatePort, idempotentPort port.IIdempotentPort) *ApplyStockUpdateService {
	return &ApplyStockUpdateService{applyStockUpdatePort, idempotentPort}
}

func (s *ApplyStockUpdateService) ApplyStockUpdate(ctx context.Context, cmd port.StockUpdateCmd) error {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       model.GoodId(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	err := s.applyStockUpdatePort.ApplyStockUpdate(goods)
	if err != nil {
		return err
	}

	if cmd.Type == port.StockUpdateCmdTypeOrder {
		idempotentCmd := port.IdempotentCmd{
			Event: "reservation",
			Id:    cmd.ReservationID,
		}
		s.idempotentPort.SaveEventID(idempotentCmd)
	}

	return nil
}
