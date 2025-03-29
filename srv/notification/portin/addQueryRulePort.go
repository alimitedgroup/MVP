package portin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/business/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/business/response"
)

type QueryRules interface {
	AddQueryRule(cmd *servicecmd.AddQueryRuleCmd) *serviceresponse.AddQueryRuleResponse
}
