package service

import (
	"github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	"github.com/alimitedgroup/MVP/srv/notification/service/response"
)

type IService interface {
	AddQueryRule(cmd *servicecmd.AddQueryRuleCmd) *serviceresponse.AddQueryRuleResponse
	AddStockUpdate(cmd *servicecmd.AddStockUpdateCmd) (*serviceresponse.AddStockUpdateResponse, error)
	GetAllQueryRules() []servicecmd.AddQueryRuleCmd
	GetCurrentQuantityByGoodID(goodID string) *serviceresponse.GetRuleResultResponse
}
