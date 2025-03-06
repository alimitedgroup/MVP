package service_Response

type SetMultipleGoodsQuantityResponse struct {
	result  string //Result of the operation: can be Success or Errors. In the first case the id slice is empty, otherwise there will be some values
	wrongID []string
}

func NewSetMultipleGoodsQuantityResponse(result string, wrongID []string) *SetMultipleGoodsQuantityResponse {
	return &SetMultipleGoodsQuantityResponse{result: result, wrongID: wrongID}
}
