package serviceresponse

type SetGoodQuantityResponse struct {
	result error
}

func NewSetGoodQuantityResponse(err error) *SetGoodQuantityResponse {
	return &SetGoodQuantityResponse{result: err}
}

func (acr *SetGoodQuantityResponse) GetOperationResult() error {
	return acr.result
}
