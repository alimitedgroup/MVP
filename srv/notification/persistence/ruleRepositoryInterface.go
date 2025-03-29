package persistence

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/types"
)

type IRuleRepository interface {
	AddRule(cmd *servicecmd.QueryRule) error
	GetAllRules() []servicecmd.QueryRule
}
