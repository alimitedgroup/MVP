package serviceresponse

type AddOrChangeResponse struct {
	result error
}

func NewAddOrChangeResponse(err error) *AddOrChangeResponse {
	return &AddOrChangeResponse{result: err}
}

func (acr *AddOrChangeResponse) GetOperationResult() error {
	return acr.result
}
