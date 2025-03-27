package serviceportin

import (
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/service/response"
)

// IAddQueryRuleUseCase rappresenta il caso d’uso “Aggiungi una regola di query”
type IAddQueryRuleUseCase interface {
	AddQueryRule(cmd *servicecmd.AddQueryRuleCmd) *serviceresponse.AddQueryRuleResponse
}
