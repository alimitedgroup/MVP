package portout

import servicecmd "github.com/alimitedgroup/MVP/srv/notification/business/cmd"

type IRuleRepository interface {
	AddRule(cmd *servicecmd.AddQueryRuleCmd) error
	GetAllRules() []servicecmd.AddQueryRuleCmd
}
