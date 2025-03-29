package serviceresponse

type AddQueryRuleResponse struct {
	result error
}

func NewAddQueryRuleResponse(err error) *AddQueryRuleResponse {
	return &AddQueryRuleResponse{result: err}
}

func (aqr *AddQueryRuleResponse) GetOperationResult() error {
	return aqr.result
}
