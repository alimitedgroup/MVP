package portout

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
)

type RuleRepository interface {
	AddRule(data types.QueryRule) (uuid.UUID, error)
	GetRule(id uuid.UUID) (types.QueryRule, error)
	ListRules() ([]types.QueryRuleWithId, error)
	EditRule(id uuid.UUID, data types.EditRule) error
	RemoveRule(id uuid.UUID) error
}
