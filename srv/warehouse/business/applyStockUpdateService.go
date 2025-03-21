package business

import (
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

func (s *ApplyStockUpdateService) ApplyStockUpdate(cmd port.StockUpdateCmd) {
	goods := make([]model.GoodStock, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.GoodStock{
			ID:       good.GoodID,
			Quantity: good.Quantity,
		})
	}

	s.applyStockUpdatePort.ApplyStockUpdate(goods)

	if cmd.Type == port.StockUpdateCmdTypeOrder || cmd.Type == port.StockUpdateCmdTypeTransfer {
		idempotentCmd := port.IdempotentCmd{
			Event: "reservation",
			Id:    cmd.ReservationID,
		}
		s.idempotentPort.SaveEventID(idempotentCmd)
	}
}
