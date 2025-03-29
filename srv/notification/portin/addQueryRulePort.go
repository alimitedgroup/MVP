package portin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/service/response"
)

type QueryRules interface {
	AddQueryRule(cmd *servicecmd.AddQueryRuleCmd) *serviceresponse.AddQueryRuleResponse
}
