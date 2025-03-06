package service

type CatalogService struct {
	addOrChangeGoodDataPort IAddOrChangeGoodDataPort
	setGoodQuantity         ISetGoodQuantityPort
	getGoodsQuantity        IGetGoodsQuantityPort
	getGoodsInfo            IGetGoodsQuantityPort
}

func NewCatalogService(AddOrChangeGoodDataPort IAddOrChangeGoodDataPort, SetGoodQuantity ISetGoodQuantityPort, GetGoodsQuantity IGetGoodsQuantityPort, GetGoodsInfo IGetGoodsQuantityPort) *CatalogService {
	return &CatalogService{addOrChangeGoodDataPort: AddOrChangeGoodDataPort, setGoodQuantity: SetGoodQuantity, getGoodsQuantity: GetGoodsQuantity, getGoodsInfo: GetGoodsInfo}
}

func (cs *CatalogService) AddOrChangeGoodData(agc *AddChangeGoodCmd) *AddOrChangeResponse {
	return cs.addOrChangeGoodDataPort.AddOrChangeGoodData(agc)
}

func checkErrinSlice(errorSlice []string) []int {
	result := []int{}
	for i := range errorSlice {
		if errorSlice[i] != "Success" {
			result = append(result, i)
		}
	}
	return result
}

func (cs *CatalogService) SetMultipleGoodsQuantity(cmd *MultipleGoodsQuantityCmd) *SetMultipleGoodsQuantityResult {
	warehouseID := cmd.GetWarehouseID()
	goodsSlice := cmd.GetGoods()
	var errorSlice []string
	var err string
	for i := range goodsSlice {
		err = cs.setGoodQuantity.SetGoodQuantity(NewSetGoodQuantityCmd(warehouseID, goodsSlice[i].GoodID, goodsSlice[i].Quantity)).GetOperationResult()
		errorSlice = append(errorSlice, err)
	}

	errors := checkErrinSlice(errorSlice)

	if len(errors) == 0 {
		return NewSetMultipleGoodsQuantityResult("Success", []string{})
	}

	var wrongID []string
	for i := range errors {
		wrongID = append(wrongID, goodsSlice[i].GoodID)
	}
	return NewSetMultipleGoodsQuantityResult("Errors", wrongID)
}
