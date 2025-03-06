package catalogAdapter

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
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
		return catalogCommon.NewCustomError("Not a valid goodID")
	}
	return nil
}

func (fr *fakeRepo) AddGood(goodID string, name string, description string) error {
	if goodID == "test-wrong-ID" {
		return catalogCommon.NewCustomError("Not a valid goodID")
	}
	return nil
}

//FINE DESCRIZIONE MOCK REPO

func TestAddChangeGoodData(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := service_Cmd.NewAddChangeGoodCmd("test-ID", "test-name", "test-desc")
			response := a.AddOrChangeGoodData(cmd)
			assert.Equal(t, response.GetOperationResult(), "Success")
		}),
	)
}

func TestAddChangeGoodDataWithWrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := service_Cmd.NewAddChangeGoodCmd("test-wrong-ID", "test-name", "test-desc")
			response := a.AddOrChangeGoodData(cmd)
			assert.Equal(t, response.GetOperationResult(), "Not a valid goodID")
		}),
	)
}

func TestSetGoodQuantity(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := service_Cmd.NewSetGoodQuantityCmd("warehouse-ID", "test-ID", 7)
			response := a.SetGoodQuantity(cmd)
			assert.Equal(t, response.GetOperationResult(), "Success")
		}),
	)
}

func TestSetGoodQuantityWithWrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := service_Cmd.NewSetGoodQuantityCmd("warehouse-ID", "test-wrong-ID", 7)
			response := a.SetGoodQuantity(cmd)
			assert.Equal(t, response.GetOperationResult(), "Not a valid goodID")
		}),
	)
}

func TestGetGoodsQuantity(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := service_Cmd.NewGetGoodsQuantityCmd()
			response := a.GetGoodsQuantity(cmd)
			assert.Equal(t, response.GetMap()["test-ID"], int64(7))
		}),
	)
}

func TestGetGoodsInfo(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeRepo,
				fx.As(new(IGoodRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogRepositoryAdapter) {
			cmd := service_Cmd.NewGetGoodsInfoCmd()
			response := a.GetGoodsInfo(cmd)
			goods := make(map[string]catalogCommon.Good)
			goods["test-ID"] = *catalogCommon.NewGood("test-ID", "test-name", "test-description")
			assert.Equal(t, response.GetMap(), goods)
		}),
	)
}
