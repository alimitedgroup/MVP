package service_Response

type SetGoodQuantityResponse struct {
	result string
}

func NewSetGoodQuantityResponse(text string) *SetGoodQuantityResponse {
	return &SetGoodQuantityResponse{result: text}
}

func (acr *SetGoodQuantityResponse) GetOperationResult() string {
	return acr.result
}
