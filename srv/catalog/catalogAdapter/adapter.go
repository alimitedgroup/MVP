package catalogAdapter

import "github.com/alimitedgroup/MVP/srv/catalog/service"

type CatalogRepositoryAdapter struct {
	repo IGoodRepository
}

func NewCatalogRepositoryAdapter(repo IGoodRepository) *CatalogRepositoryAdapter {
	return &CatalogRepositoryAdapter{repo: repo}
}

func (cra *CatalogRepositoryAdapter) AddOrChangeGoodData(agc *service.AddChangeGoodCmd) *service.AddOrChangeResponse {
	err := cra.repo.AddGood(agc.GetId(), agc.GetName(), agc.GetDescription())
	if err != nil {
		return service.NewAddOrChangeResponse(err.Error())
	}
	return service.NewAddOrChangeResponse("Success")
}

func (cra *CatalogRepositoryAdapter) SetGoodQuantity(agqc *service.SetGoodQuantityCmd) *service.SetGoodQuantityResponse {
	err := cra.repo.SetGoodQuantity(agqc.GetWarehouseId(), agqc.GetGoodId(), agqc.GetNewQuantity())
	if err != nil {
		return service.NewSetGoodQuantityResponse(err.Error())
	}
	return service.NewSetGoodQuantityResponse("Success")
}

func (cra *CatalogRepositoryAdapter) GetGoodsQuantity(ggqc *service.GetGoodsQuantityCmd) *service.GetGoodsQuantityResponse {
	return service.NewGetGoodsQuantityResponse(cra.repo.GetGoodsGlobalQuantity())
}

func (cra *CatalogRepositoryAdapter) GetGoodsInfo(ggqc *service.GetGoodsInfoCmd) *service.GetGoodsInfoResponse {
	return service.NewGetGoodsInfoResponse(cra.repo.GetGoods())
}
