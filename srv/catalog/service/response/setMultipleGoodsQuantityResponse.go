package serviceresponse

type SetMultipleGoodsQuantityResponse struct {
	result  error //Result of the operation: can be Success or Errors. In the first case the id slice is empty, otherwise there will be some values
	wrongID []string
}

func NewSetMultipleGoodsQuantityResponse(err error, wrongID []string) *SetMultipleGoodsQuantityResponse {
	return &SetMultipleGoodsQuantityResponse{result: err, wrongID: wrongID}
}

func (smgqr *SetMultipleGoodsQuantityResponse) GetOperationResult() error {
	return smgqr.result
}

func (smgqr *SetMultipleGoodsQuantityResponse) GetWrongIDSlice() []string {
	return smgqr.wrongID
}
