package service

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/notification/service/portin"
	serviceportout "github.com/alimitedgroup/MVP/srv/notification/service/portout"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/service/response"
)

type NotificationService struct {
	ruleRepo       serviceportout.IRuleRepository
	alertPublisher serviceportout.IStockEventPublisher
	quantityReader serviceportout.IRuleQueryRepository
	stockRepo      serviceportout.IStockRepository
}

func NewNotificationService(
	ruleRepo serviceportout.IRuleRepository,
	alertPublisher serviceportout.IStockEventPublisher,
	quantityReader serviceportout.IRuleQueryRepository,
	stockRepo serviceportout.IStockRepository,
) *NotificationService {
	return &NotificationService{
		ruleRepo:       ruleRepo,
		alertPublisher: alertPublisher,
		quantityReader: quantityReader,
		stockRepo:      stockRepo,
	}
}

var _ serviceportin.IAddQueryRuleUseCase = (*NotificationService)(nil)
var _ serviceportin.IAddStockUpdateUseCase = (*NotificationService)(nil)

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

func (ns *NotificationService) PublishStockAlert(alert serviceportout.StockAlertEvent) error {
	return ns.alertPublisher.PublishStockAlert(alert)
}
