package portout

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
)

type IRuleQueryRepository interface {
	GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse
}
