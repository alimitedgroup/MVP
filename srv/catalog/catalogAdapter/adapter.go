package catalogAdapter

import (
	"github.com/alimitedgroup/MVP/srv/catalog/persistence"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type CatalogRepositoryAdapter struct {
	repo persistence.IGoodRepository
}

func NewCatalogRepositoryAdapter(repo persistence.IGoodRepository) *CatalogRepositoryAdapter {
	return &CatalogRepositoryAdapter{repo: repo}
}

func (cra *CatalogRepositoryAdapter) SetGoodQuantity(agqc *servicecmd.SetGoodQuantityCmd) *serviceresponse.SetGoodQuantityResponse {
	err := cra.repo.SetGoodQuantity(agqc.GetWarehouseId(), agqc.GetGoodId(), agqc.GetNewQuantity())
	if err != nil {
		return serviceresponse.NewSetGoodQuantityResponse(err)
	}
	return serviceresponse.NewSetGoodQuantityResponse(nil)
}

func (cra *CatalogRepositoryAdapter) GetGoodsQuantity(ggqc *servicecmd.GetGoodsQuantityCmd) *serviceresponse.GetGoodsQuantityResponse {
	return serviceresponse.NewGetGoodsQuantityResponse(cra.repo.GetGoodsGlobalQuantity())
}

func (cra *CatalogRepositoryAdapter) GetWarehouses(*servicecmd.GetWarehousesCmd) *serviceresponse.GetWarehousesResponse {
	return serviceresponse.NewGetWarehousesResponse(cra.repo.GetWarehouses())
}
