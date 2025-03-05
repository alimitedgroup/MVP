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
