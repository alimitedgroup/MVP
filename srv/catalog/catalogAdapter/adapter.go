package catalogAdapter

import "github.com/alimitedgroup/MVP/srv/catalog/service"

type catalogRepositoryAdapter struct {
	repo IGoodRepository
}

func NewCatalogRepositoryAdapter(repo IGoodRepository) *catalogRepositoryAdapter {
	return &catalogRepositoryAdapter{repo: repo}
}

func (cra *catalogRepositoryAdapter) AddOrChangeGoodData(agc *service.AddGoodCmd) *service.AddOrChangeResponse {
	err := cra.repo.AddGood(agc.GetId(), agc.GetName(), agc.GetDescription())
	if err != nil {
		return service.NewAddOrChangeResponse(err.Error())
	}
	return service.NewAddOrChangeResponse("Success")
}

func (cra *catalogRepositoryAdapter) SetGoodQuantity(agqc *service.SetGoodQuantityCmd) *service.SetGoodQuantityResponse {
	err := cra.repo.SetGoodQuantity(agqc.GetWarehouseId(), agqc.GetGoodId(), agqc.GetNewQuantity())
	if err != nil {
		return service.NewSetGoodQuantityResponse(err.Error())
	}
	return service.NewSetGoodQuantityResponse("Success")
}

func (cra *catalogRepositoryAdapter) GetGoodsQuantity(ggqc *service.GetGoodsQuantityCmd) *service.GetGoodsQuantityResponse {
	return service.NewGetGoodsQuantityResponse(cra.repo.GetGoodsGlobalQuantity())
}

func (cra *catalogRepositoryAdapter) GetGoodsInfo(ggqc *service.GetGoodsQuantityCmd) *service.GetGoodsInfoResponse {
	return service.NewGetGoodsInfoResponse(cra.repo.GetGoods())
}
