package serviceresponse

type AddStockUpdateResponse struct {
	result error
}

func NewAddStockUpdateResponse(err error) *AddStockUpdateResponse {
	return &AddStockUpdateResponse{result: err}
}

func (asr *AddStockUpdateResponse) GetOperationResult() error {
	return asr.result
}
