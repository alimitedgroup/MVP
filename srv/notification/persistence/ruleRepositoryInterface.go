package persistence

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
)

type IRuleRepository interface {
	AddRule(cmd types.QueryRule) (uuid.UUID, error)
	ListRules() ([]types.QueryRuleWithId, error)
}
