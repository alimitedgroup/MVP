package persistence

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func TestAddGood(t *testing.T) {
	// Testa l'aggiunta di una merce
	fx.New(
		fx.Provide(NewCatalogGoodDataRepository),
		fx.Invoke(func(cr *CatalogGoodDataRepository) {
			err2 := cr.AddGood("test-ID", "test-name", "test-description")
			resultGoods := cr.GetGoods()
			assert.Equal(t, err2, nil)
			assert.Equal(t, resultGoods["test-ID"].GetID(), "test-ID")
			assert.Equal(t, resultGoods["test-ID"].GetName(), "test-name")
			assert.Equal(t, resultGoods["test-ID"].GetDescription(), "test-description")
		}),
	)
}

func TestChangeGoodData(t *testing.T) {
	// Testa la modifica di una merce
	fx.New(
		fx.Provide(NewCatalogGoodDataRepository),
		fx.Invoke(func(cr *CatalogGoodDataRepository) {
			err1 := cr.AddGood("test-ID", "test-name", "test-description")
			err3 := cr.AddGood("test-ID", "new-test-name", "new-test-description")
			resultGoods := cr.GetGoods()
			assert.Equal(t, err1, nil)
			assert.Equal(t, err3, nil)
			assert.Equal(t, resultGoods["test-ID"].GetID(), "test-ID")
			assert.Equal(t, resultGoods["test-ID"].GetName(), "new-test-name")
			assert.Equal(t, resultGoods["test-ID"].GetDescription(), "new-test-description")
		}),
	)
}

func TestChangeGoodDataWrongID(t *testing.T) {
	// Testa la modifica di una merce
	fx.New(
		fx.Provide(NewCatalogGoodDataRepository),
		fx.Invoke(func(cr *CatalogGoodDataRepository) {
			err1 := cr.AddGood("test-ID", "test-name", "test-description")
			err2 := cr.changeGoodData("2test-ID", "test-name", "test-description")
			assert.Equal(t, err1, nil)
			assert.Equal(t, err2, catalogCommon.ErrGoodIdNotValid)
		}),
	)
}

func TestChangeGoodDataEmptyName(t *testing.T) {
	// Testa la modifica di una merce
	fx.New(
		fx.Provide(NewCatalogGoodDataRepository),
		fx.Invoke(func(cr *CatalogGoodDataRepository) {
			err1 := cr.AddGood("test-ID", "test-name", "test-description")
			err2 := cr.changeGoodData("test-ID", "", "test-description")
			assert.Equal(t, err1, nil)
			assert.Equal(t, err2, dto.ErrEmptyName)
		}),
	)
}

func TestChangeGoodDataEmptyDescription(t *testing.T) {
	// Testa la modifica di una merce
	fx.New(
		fx.Provide(NewCatalogGoodDataRepository),
		fx.Invoke(func(cr *CatalogGoodDataRepository) {
			err1 := cr.AddGood("test-ID", "test-name", "test-description")
			err2 := cr.changeGoodData("test-ID", "2test-name", "")
			assert.Equal(t, err1, nil)
			assert.Equal(t, err2, dto.ErrEmptyDescription)
		}),
	)
}
