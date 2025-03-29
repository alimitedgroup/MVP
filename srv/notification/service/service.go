package service

import (
	"github.com/alimitedgroup/MVP/srv/notification/portin"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	"github.com/alimitedgroup/MVP/srv/notification/service/response"
)

type NotificationService struct {
	ruleRepo       portout.IRuleRepository
	alertPublisher portout.IStockEventPublisher
	quantityReader portout.IRuleQueryRepository
	stockRepo      portout.IStockRepository
}

func NewNotificationService(
	ruleRepo portout.IRuleRepository,
	alertPublisher portout.IStockEventPublisher,
	quantityReader portout.IRuleQueryRepository,
	stockRepo portout.IStockRepository,
) *NotificationService {
	return &NotificationService{
		ruleRepo:       ruleRepo,
		alertPublisher: alertPublisher,
		quantityReader: quantityReader,
		stockRepo:      stockRepo,
	}
}

var _ portin.QueryRules = (*NotificationService)(nil)
var _ portin.StockUpdates = (*NotificationService)(nil)

func (ns *NotificationService) AddQueryRule(cmd *servicecmd.AddQueryRuleCmd) *serviceresponse.AddQueryRuleResponse {
	err := ns.ruleRepo.AddRule(cmd)
	return serviceresponse.NewAddQueryRuleResponse(err)
}

func (ns *NotificationService) AddStockUpdate(cmd *servicecmd.AddStockUpdateCmd) (*serviceresponse.AddStockUpdateResponse, error) {
	return ns.stockRepo.SaveStockUpdate(cmd), nil
}

// ========== Utility per RuleChecker ==========

func (ns *NotificationService) GetAllQueryRules() []servicecmd.AddQueryRuleCmd {
	return ns.ruleRepo.GetAllRules()
}

func (ns *NotificationService) GetCurrentQuantityByGoodID(goodID string) *serviceresponse.GetRuleResultResponse {
	return ns.quantityReader.GetCurrentQuantityByGoodID(goodID)
}

func (ns *NotificationService) PublishStockAlert(alert portout.StockAlertEvent) error {
	return ns.alertPublisher.PublishStockAlert(alert)
}
