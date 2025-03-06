package service

type SetMultipleGoodsQuantityResult struct {
	result  string //Result of the operation: can be Success or Errors. In the first case the id slice is empty, otherwise there will be some values
	wrongID []string
}

func NewSetMultipleGoodsQuantityResult(result string, wrongID []string) *SetMultipleGoodsQuantityResult {
	return &SetMultipleGoodsQuantityResult{result: result, wrongID: wrongID}
}
