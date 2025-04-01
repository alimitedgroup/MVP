package portin

import (
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
)

type QueryRules interface {
	AddQueryRule(cmd types.QueryRule) (uuid.UUID, error)
	GetQueryRule(id uuid.UUID) (types.QueryRule, error)
	ListQueryRules() ([]types.QueryRuleWithId, error)
	EditQueryRule(id uuid.UUID, data types.EditRule) error
	RemoveQueryRule(id uuid.UUID) error
}
