package business

import (
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
)

type Business struct {
	ruleRepo       portout.IRuleRepository
	alertPublisher portout.IStockEventPublisher
	quantityReader portout.IRuleQueryRepository
	stockRepo      portout.IStockRepository
}

func NewBusiness(
	ruleRepo portout.IRuleRepository,
	alertPublisher portout.IStockEventPublisher,
	quantityReader portout.IRuleQueryRepository,
	stockRepo portout.IStockRepository,
) *Business {
	return &Business{
		ruleRepo:       ruleRepo,
		alertPublisher: alertPublisher,
		quantityReader: quantityReader,
		stockRepo:      stockRepo,
	}
}

// Asserzione a compile-time che Business implementi le interfacce delle port-in
var _ portin.QueryRules = (*Business)(nil)
var _ portin.StockUpdates = (*Business)(nil)

func (ns *Business) AddQueryRule(cmd *types.QueryRule) *types.AddQueryRuleResponse {
	err := ns.ruleRepo.AddRule(cmd)
	return types.NewAddQueryRuleResponse(err)
}

func (ns *Business) AddStockUpdate(cmd *types.AddStockUpdateCmd) (*types.AddStockUpdateResponse, error) {
	return ns.stockRepo.SaveStockUpdate(cmd), nil
}

// ========== Utility per RuleChecker ==========

func (ns *Business) GetAllQueryRules() []types.QueryRule {
	return ns.ruleRepo.GetAllRules()
}

func (ns *Business) GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse {
	return ns.quantityReader.GetCurrentQuantityByGoodID(goodID)
}

func (ns *Business) PublishStockAlert(alert types.StockAlertEvent) error {
	return ns.alertPublisher.PublishStockAlert(alert)
}
