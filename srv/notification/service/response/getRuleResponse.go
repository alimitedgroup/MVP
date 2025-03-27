package serviceresponse

type GetRuleResultResponse struct {
	GoodID          string
	CurrentQuantity int
	Err             error
}

func NewGetRuleResultResponse(goodID string, quantity int, err error) *GetRuleResultResponse {
	return &GetRuleResultResponse{
		GoodID:          goodID,
		CurrentQuantity: quantity,
		Err:             err,
	}
}
