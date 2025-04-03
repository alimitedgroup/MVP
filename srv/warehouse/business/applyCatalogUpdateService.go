package business

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
)

type ApplyCatalogUpdateService struct {
	applyCatalogUpdatePort port.IApplyCatalogUpdatePort
	transactionPort        port.ITransactionPort
}

func NewApplyCatalogUpdateService(applyGoodUpdatePort port.IApplyCatalogUpdatePort, transactionPort port.ITransactionPort) *ApplyCatalogUpdateService {
	return &ApplyCatalogUpdateService{applyGoodUpdatePort, transactionPort}
}

func (s *ApplyCatalogUpdateService) ApplyCatalogUpdate(cmd port.CatalogUpdateCmd) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	good := model.GoodInfo{
		ID:          cmd.GoodID,
		Name:        cmd.Name,
		Description: cmd.Description,
	}
	s.applyCatalogUpdatePort.ApplyCatalogUpdate(good)
}
