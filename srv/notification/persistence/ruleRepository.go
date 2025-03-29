package persistence

import (
	"sync"

	servicecmd "github.com/alimitedgroup/MVP/srv/notification/business/cmd"
)

type RuleRepository struct {
	ruleMap map[string]servicecmd.AddQueryRuleCmd
	mutex   sync.RWMutex
}

func NewRuleRepository() *RuleRepository {
	return &RuleRepository{
		ruleMap: make(map[string]servicecmd.AddQueryRuleCmd),
	}
}

func (rr *RuleRepository) AddRule(cmd *servicecmd.AddQueryRuleCmd) error {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	rr.ruleMap[cmd.GetGoodID()] = *cmd
	return nil
}

func (rr *RuleRepository) GetAllRules() []servicecmd.AddQueryRuleCmd {
	rr.mutex.RLock()
	defer rr.mutex.RUnlock()

	rules := make([]servicecmd.AddQueryRuleCmd, 0, len(rr.ruleMap))
	for _, rule := range rr.ruleMap {
		rules = append(rules, rule)
	}
	return rules
}
