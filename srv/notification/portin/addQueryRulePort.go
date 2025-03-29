package portin

import (
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/types"
)

type QueryRules interface {
	AddQueryRule(cmd *serviceresponse.QueryRule) *serviceresponse.AddQueryRuleResponse
}
