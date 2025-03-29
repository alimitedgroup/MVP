package portout

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
)

type IRuleRepository interface {
	AddRule(cmd *types.QueryRule) error
	GetAllRules() []types.QueryRule
}
