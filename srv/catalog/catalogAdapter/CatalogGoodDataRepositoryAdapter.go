package catalogAdapter

import (
	"github.com/alimitedgroup/MVP/srv/catalog/persistence"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
)

type CatalogGoodDataRepositoryAdapter struct {
	repo persistence.ICatalogGoodDataRepository
}

func NewCatalogGoodDataRepositoryAdapter(repo persistence.ICatalogGoodDataRepository) *CatalogGoodDataRepositoryAdapter {
	return &CatalogGoodDataRepositoryAdapter{repo: repo}
}

func (cra *CatalogGoodDataRepositoryAdapter) AddOrChangeGoodData(agc *servicecmd.AddChangeGoodCmd) *serviceresponse.AddOrChangeResponse {
	err := cra.repo.AddGood(agc.GetId(), agc.GetName(), agc.GetDescription())
	if err != nil {
		return serviceresponse.NewAddOrChangeResponse(err)
	}
	return serviceresponse.NewAddOrChangeResponse(nil)
}

func (cra *CatalogGoodDataRepositoryAdapter) GetGoodsInfo(ggqc *servicecmd.GetGoodsInfoCmd) *serviceresponse.GetGoodsInfoResponse {
	return serviceresponse.NewGetGoodsInfoResponse(cra.repo.GetGoods())
}
