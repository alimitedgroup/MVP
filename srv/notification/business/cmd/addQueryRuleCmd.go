package servicecmd

type AddQueryRuleCmd struct {
	goodID    string
	operator  string
	threshold int
}

func NewAddQueryRuleCmd(goodID string, operator string, threshold int) *AddQueryRuleCmd {
	return &AddQueryRuleCmd{
		goodID:    goodID,
		operator:  operator,
		threshold: threshold,
	}
}

func (aqr *AddQueryRuleCmd) GetGoodID() string {
	return aqr.goodID
}

func (aqr *AddQueryRuleCmd) GetOperator() string {
	return aqr.operator
}

func (aqr *AddQueryRuleCmd) GetThreshold() int {
	return aqr.threshold
}
