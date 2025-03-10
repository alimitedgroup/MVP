package catalogAdapter

import (
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
)

type CatalogRepositoryAdapter struct {
	repo IGoodRepository
}

func NewCatalogRepositoryAdapter(repo IGoodRepository) *CatalogRepositoryAdapter {
	return &CatalogRepositoryAdapter{repo: repo}
}

func (cra *CatalogRepositoryAdapter) AddOrChangeGoodData(agc *service_Cmd.AddChangeGoodCmd) *service_Response.AddOrChangeResponse {
	err := cra.repo.AddGood(agc.GetId(), agc.GetName(), agc.GetDescription())
	if err != nil {
		return service_Response.NewAddOrChangeResponse(err.Error())
	}
	return service_Response.NewAddOrChangeResponse("Success")
}

func (cra *CatalogRepositoryAdapter) SetGoodQuantity(agqc *service_Cmd.SetGoodQuantityCmd) *service_Response.SetGoodQuantityResponse {
	err := cra.repo.SetGoodQuantity(agqc.GetWarehouseId(), agqc.GetGoodId(), agqc.GetNewQuantity())
	if err != nil {
		return service_Response.NewSetGoodQuantityResponse(err.Error())
	}
	return service_Response.NewSetGoodQuantityResponse("Success")
}

func (cra *CatalogRepositoryAdapter) GetGoodsQuantity(ggqc *service_Cmd.GetGoodsQuantityCmd) *service_Response.GetGoodsQuantityResponse {
	return service_Response.NewGetGoodsQuantityResponse(cra.repo.GetGoodsGlobalQuantity())
}

func (cra *CatalogRepositoryAdapter) GetGoodsInfo(ggqc *service_Cmd.GetGoodsInfoCmd) *service_Response.GetGoodsInfoResponse {
	return service_Response.NewGetGoodsInfoResponse(cra.repo.GetGoods())
}

func (cra *CatalogRepositoryAdapter) GetWarehouses(*service_Cmd.GetWarehousesCmd) *service_Response.GetWarehousesResponse {
	return service_Response.NewGetWarehousesResponse(cra.repo.GetWarehouses())
}
