package portout

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
)

type RuleQueryRepository interface {
	GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse
}
