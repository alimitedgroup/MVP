package persistence

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/types"
	"sync"
)

type RuleRepository struct {
	ruleMap map[string]servicecmd.QueryRule
	mutex   sync.RWMutex
}

func NewRuleRepository() *RuleRepository {
	return &RuleRepository{
		ruleMap: make(map[string]servicecmd.QueryRule),
	}
}

func (rr *RuleRepository) AddRule(cmd *servicecmd.QueryRule) error {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	rr.ruleMap[cmd.GoodId] = *cmd
	return nil
}

func (rr *RuleRepository) GetAllRules() []servicecmd.QueryRule {
	rr.mutex.RLock()
	defer rr.mutex.RUnlock()

	rules := make([]servicecmd.QueryRule, 0, len(rr.ruleMap))
	for _, rule := range rr.ruleMap {
		rules = append(rules, rule)
	}
	return rules
}
