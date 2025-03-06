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

func (cs *CatalogService) SetMultipleGoodsQuantity(cmd *MultipleGoodsQuantityCmd) {
	warehouseID := cmd.GetWarehouseID()
	goodsSlice := cmd.GetGoods()

	for i := range goodsSlice {
		cs.setGoodQuantity.SetGoodQuantity(NewSetGoodQuantityCmd(warehouseID, goodsSlice[i].GoodID, goodsSlice[i].Quantity))
	}
}
