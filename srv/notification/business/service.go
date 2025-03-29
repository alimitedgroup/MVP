package business

import (
	"github.com/alimitedgroup/MVP/srv/notification/business/cmd"
	"github.com/alimitedgroup/MVP/srv/notification/business/response"
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
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

func (ns *Business) AddQueryRule(cmd *servicecmd.AddQueryRuleCmd) *serviceresponse.AddQueryRuleResponse {
	err := ns.ruleRepo.AddRule(cmd)
	return serviceresponse.NewAddQueryRuleResponse(err)
}

func (ns *Business) AddStockUpdate(cmd *servicecmd.AddStockUpdateCmd) (*serviceresponse.AddStockUpdateResponse, error) {
	return ns.stockRepo.SaveStockUpdate(cmd), nil
}

// ========== Utility per RuleChecker ==========

func (ns *Business) GetAllQueryRules() []servicecmd.AddQueryRuleCmd {
	return ns.ruleRepo.GetAllRules()
}

func (ns *Business) GetCurrentQuantityByGoodID(goodID string) *serviceresponse.GetRuleResultResponse {
	return ns.quantityReader.GetCurrentQuantityByGoodID(goodID)
}

func (ns *Business) PublishStockAlert(alert portout.StockAlertEvent) error {
	return ns.alertPublisher.PublishStockAlert(alert)
}
