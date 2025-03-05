package service

type AddOrChangeResponse struct {
	result string
}

func NewAddOrChangeResponse(text string) *AddOrChangeResponse {
	return &AddOrChangeResponse{result: text}
}

func (acr *AddOrChangeResponse) GetOperationResult() string {
	return acr.result
}
