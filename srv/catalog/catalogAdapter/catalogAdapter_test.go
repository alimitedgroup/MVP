package catalogAdapter

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"

	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	persistence "github.com/alimitedgroup/MVP/srv/catalog/persistence"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

// INIZIO DESCRIZIONE MOCK REPO
type fakeRepo struct {
	//simula un repo funzionante
}

func NewFakeRepo() *fakeRepo {
	return &fakeRepo{}
}

func (fr *fakeRepo) GetGoodsGlobalQuantity() map[string]int64 {
	goods := make(map[string]int64)
	goods["test-ID"] = int64(7)
	return goods
}

func (fr *fakeRepo) SetGoodQuantity(warehouseID string, goodID string, newQuantity int64) error {
	if goodID == "test-wrong-ID" {
		return catalogCommon.ErrGoodIdNotValid
	}
	return nil
}

func (fr *fakeRepo) GetWarehouses() map[string]dto.Warehouse {
	warehouses := make(map[string]dto.Warehouse)
	warehouses["test-warehouse-ID"] = *dto.NewWarehouse("test-warehose-ID")
	return warehouses
}

//FINE DESCRIZIONE MOCK REPO

func TestGetWarehouses(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(persistence.IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := servicecmd.NewGetWarehousesCmd()
			response := a.GetWarehouses(cmd)
			warehouses := make(map[string]dto.Warehouse)
			warehouses["test-warehouse-ID"] = *dto.NewWarehouse("test-warehose-ID")
			assert.Equal(t, response.GetWarehouseMap(), warehouses)
		}),
	)
}

func TestSetGoodQuantity(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(persistence.IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := servicecmd.NewSetGoodQuantityCmd("warehouse-ID", "test-ID", 7)
			response := a.SetGoodQuantity(cmd)
			assert.Equal(t, response.GetOperationResult(), nil)
		}),
	)
}

func TestSetGoodQuantityWithWrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(persistence.IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := servicecmd.NewSetGoodQuantityCmd("warehouse-ID", "test-wrong-ID", 7)
			response := a.SetGoodQuantity(cmd)
			assert.Equal(t, response.GetOperationResult(), catalogCommon.ErrGoodIdNotValid)
		}),
	)
}

func TestGetGoodsQuantity(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(persistence.IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := servicecmd.NewGetGoodsQuantityCmd()
			response := a.GetGoodsQuantity(cmd)
			assert.Equal(t, response.GetMap()["test-ID"], int64(7))
		}),
	)
}
