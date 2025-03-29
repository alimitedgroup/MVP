package business

import (
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
)

func NewBusiness(
	ruleRepo portout.RuleRepository,
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

type Business struct {
	ruleRepo       portout.RuleRepository
	alertPublisher portout.IStockEventPublisher
	quantityReader portout.IRuleQueryRepository
	stockRepo      portout.IStockRepository
}

// =========== QueryRules port-in ===========

var _ portin.QueryRules = (*Business)(nil)

func (ns *Business) AddQueryRule(cmd types.QueryRule) (uuid.UUID, error) {
	return ns.ruleRepo.AddRule(cmd)
}

func (ns *Business) GetQueryRule(id uuid.UUID) (types.QueryRule, error) {
	return ns.ruleRepo.GetRule(id)
}

func (ns *Business) ListQueryRules() ([]types.QueryRuleWithId, error) {
	return ns.ruleRepo.ListRules()
}

func (ns *Business) EditQueryRule(id uuid.UUID, data types.EditRule) error {
	return ns.ruleRepo.EditRule(id, data)
}

func (ns *Business) RemoveQueryRule(id uuid.UUID) error {
	return ns.ruleRepo.RemoveRule(id)
}

// ========== StockUpdates port-in ==========

var _ portin.StockUpdates = (*Business)(nil)

func (ns *Business) RecordStockUpdate(cmd *types.AddStockUpdateCmd) error {
	return ns.stockRepo.SaveStockUpdate(cmd)
}

// ========== Utility per RuleChecker ==========

func (ns *Business) GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse {
	return ns.quantityReader.GetCurrentQuantityByGoodID(goodID)
}

func (ns *Business) PublishStockAlert(alert types.StockAlertEvent) error {
	return ns.alertPublisher.PublishStockAlert(alert)
}
