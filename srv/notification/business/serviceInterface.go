package business

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
)

type IService interface {
	AddQueryRule(cmd *types.QueryRule) error
	AddStockUpdate(cmd *types.AddStockUpdateCmd) (*types.AddStockUpdateResponse, error)
	GetAllQueryRules() []types.QueryRule
	GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse
}
