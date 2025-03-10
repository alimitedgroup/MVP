package catalogAdapter

import (
	"testing"

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

func (fr *fakeRepo) GetGoods() map[string]catalogCommon.Good {
	goods := make(map[string]catalogCommon.Good)
	goods["test-ID"] = *catalogCommon.NewGood("test-ID", "test-name", "test-description")
	return goods
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

func (fr *fakeRepo) AddGood(goodID string, name string, description string) error {
	if goodID == "test-wrong-ID" {
		return catalogCommon.ErrGoodIdNotValid
	}
	return nil
}

func (fr *fakeRepo) GetWarehouses() map[string]catalogCommon.Warehouse {
	warehouses := make(map[string]catalogCommon.Warehouse)
	warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")
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
			warehouses := make(map[string]catalogCommon.Warehouse)
			warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")
			assert.Equal(t, response.GetWarehouseMap(), warehouses)
		}),
	)
}

func TestAddChangeGoodData(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(persistence.IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := servicecmd.NewAddChangeGoodCmd("test-ID", "test-name", "test-desc")
			response := a.AddOrChangeGoodData(cmd)
			assert.Equal(t, response.GetOperationResult(), nil)
		}),
	)
}

func TestAddChangeGoodDataWithWrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(persistence.IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := servicecmd.NewAddChangeGoodCmd("test-wrong-ID", "test-name", "test-desc")
			response := a.AddOrChangeGoodData(cmd)
			assert.Equal(t, response.GetOperationResult(), catalogCommon.ErrGoodIdNotValid)
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

func TestGetGoodsInfo(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(persistence.IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := servicecmd.NewGetGoodsInfoCmd()
			response := a.GetGoodsInfo(cmd)
			goods := make(map[string]catalogCommon.Good)
			goods["test-ID"] = *catalogCommon.NewGood("test-ID", "test-name", "test-description")
			assert.Equal(t, response.GetMap(), goods)
		}),
	)
}
