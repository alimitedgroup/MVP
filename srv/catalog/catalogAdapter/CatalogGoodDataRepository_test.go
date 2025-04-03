package catalogAdapter

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	"github.com/alimitedgroup/MVP/srv/catalog/persistence"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

// INIZIO DESCRIZIONE MOCK REPO
type fakeGoodRepo struct {
	//simula un repo funzionante
}

func NewFakeGoodRepo() *fakeGoodRepo {
	return &fakeGoodRepo{}
}

func (fr *fakeRepo) GetGoods() map[string]dto.Good {
	goods := make(map[string]dto.Good)
	goods["test-ID"] = *dto.NewGood("test-ID", "test-name", "test-description")
	return goods
}

func (fr *fakeRepo) AddGood(goodID string, name string, description string) error {
	if goodID == "test-wrong-ID" {
		return catalogCommon.ErrGoodIdNotValid
	}
	return nil
}

//FINE DESCRIZIONE MOCK REPO

func TestAddChangeGoodData(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeGoodRepo,
				fx.As(new(persistence.ICatalogGoodDataRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogGoodDataRepositoryAdapter) {
			cmd := servicecmd.NewAddChangeGoodCmd("test-ID", "test-name", "test-desc")
			response := a.AddOrChangeGoodData(cmd)
			assert.Equal(t, response.GetOperationResult(), nil)
		}),
	)
}

func TestAddChangeGoodDataWithWrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeGoodRepo,
				fx.As(new(persistence.ICatalogGoodDataRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogGoodDataRepositoryAdapter) {
			cmd := servicecmd.NewAddChangeGoodCmd("test-wrong-ID", "test-name", "test-desc")
			response := a.AddOrChangeGoodData(cmd)
			assert.Equal(t, response.GetOperationResult(), catalogCommon.ErrGoodIdNotValid)
		}),
	)
}

func TestGetGoodsInfo(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeGoodRepo,
				fx.As(new(persistence.ICatalogGoodDataRepository))),
		),
		fx.Provide(NewCatalogRepositoryAdapter),
		fx.Invoke(func(a *CatalogGoodDataRepositoryAdapter) {
			cmd := servicecmd.NewGetGoodsInfoCmd()
			response := a.GetGoodsInfo(cmd)
			goods := make(map[string]dto.Good)
			goods["test-ID"] = *dto.NewGood("test-ID", "test-name", "test-description")
			assert.Equal(t, response.GetMap(), goods)
		}),
	)
}
