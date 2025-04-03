package business

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
)

type ApplyStockUpdateService struct {
	applyStockUpdatePort port.IApplyStockUpdatePort
	idempotentPort       port.IIdempotentPort
	transactionPort      port.TransactionPort
}

func NewApplyStockUpdateService(applyStockUpdatePort port.IApplyStockUpdatePort, idempotentPort port.IIdempotentPort, transactionPort port.TransactionPort) *ApplyStockUpdateService {
	return &ApplyStockUpdateService{applyStockUpdatePort, idempotentPort, transactionPort}
}

func (s *ApplyStockUpdateService) ApplyStockUpdate(cmd port.StockUpdateCmd) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

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
			ID:    cmd.ReservationID,
		}
		s.idempotentPort.SaveEventID(idempotentCmd)
	}
}
