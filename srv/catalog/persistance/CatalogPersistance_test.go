package goodRepository

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func TestAddGood(t *testing.T) {
	// Testa l'aggiunta di una merce
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *catalogRepository) {
			cr.AddGood("test-name", "test-description", "test-ID")
			cr.SetGoodQuantity("test-warehouse-ID", "test-ID", 7)
			resultGoods := cr.GetGoods()
			assert.Equal(t, resultGoods["test-ID"].GetID(), "test-ID")
			assert.Equal(t, resultGoods["test-ID"].GetName(), "test-name")
			assert.Equal(t, resultGoods["test-ID"].GetDescription(), "test-description")
			assert.Equal(t, resultGoods["test-ID"].GetGlobalQuantity(), int64(7))
		}),
	)
}

func TestChangeGoodData(t *testing.T) {
	// Testa l'aggiunta di una merce
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *catalogRepository) {
			cr.AddGood("test-name", "test-description", "test-ID")
			cr.SetGoodQuantity("test-warehouse-ID", "test-ID", 7)
			resultGoods := cr.GetGoods()

			assert.Equal(t, resultGoods["test-ID"].GetID(), "test-ID")
			assert.Equal(t, resultGoods["test-ID"].GetName(), "test-name")
			assert.Equal(t, resultGoods["test-ID"].GetDescription(), "test-description")
			assert.Equal(t, resultGoods["test-ID"].GetGlobalQuantity(), int64(7))
		}),
	)
}

func TestAddWarehouse(t *testing.T) {
	//testa se l'arrivo di un messaggio con un warehouse non memorizzato determina l'aggiunta del warehouse
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *catalogRepository) {
			cr.AddGood("test-name", "test-description", "test-ID")
			cr.SetGoodQuantity("test-warehouse-ID", "test-name", 7)
			_, presence := cr.GetWarehouses()["test-warehouse-ID"]
			assert.Equal(t, presence, true)
		}),
	)
}
