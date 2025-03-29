package portout

import servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"

type IRuleRepository interface {
	AddRule(cmd *servicecmd.AddQueryRuleCmd) error
	GetAllRules() []servicecmd.AddQueryRuleCmd
}
