package business

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
)

type IService interface {
	AddQueryRule(cmd *types.AddQueryRuleCmd) *types.AddQueryRuleResponse
	AddStockUpdate(cmd *types.AddStockUpdateCmd) (*types.AddStockUpdateResponse, error)
	GetAllQueryRules() []types.AddQueryRuleCmd
	GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse
}
