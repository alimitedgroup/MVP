package persistence

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func TestGetGoodsGlobalQt(t *testing.T) {
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *CatalogRepository) {
			err4 := cr.SetGoodQuantity("mag1", "test-ID", 7)
			err5 := cr.SetGoodQuantity("mag1", "2test-ID", 8)
			err6 := cr.SetGoodQuantity("mag2", "test-ID", 9)
			err7 := cr.SetGoodQuantity("mag2", "2test-ID", 2)
			err8 := cr.SetGoodQuantity("mag3", "3test-ID", 7)
			err9 := cr.SetGoodQuantity("mag3", "3test-ID", 3)
			result := cr.GetGoodsGlobalQuantity()
			assert.Equal(t, result["test-ID"], int64(16))
			assert.Equal(t, result["2test-ID"], int64(10))
			assert.Equal(t, result["3test-ID"], int64(3))
			assert.Equal(t, err4, nil)
			assert.Equal(t, err5, nil)
			assert.Equal(t, err6, nil)
			assert.Equal(t, err7, nil)
			assert.Equal(t, err8, nil)
			assert.Equal(t, err9, nil)
		}),
	)
}

func TestAddGoodQuantity(t *testing.T) {
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *CatalogRepository) {
			err := cr.SetGoodQuantity("test-warehouse-ID", "test-ID", 7)
			assert.Equal(t, err, nil)
			assert.Equal(t, cr.GetGoodsGlobalQuantity()["test-ID"], int64(7))
		},
		),
	)
}

func TestAddWarehouse(t *testing.T) {
	//testa se l'arrivo di un messaggio con un warehouse non memorizzato determina l'aggiunta del warehouse
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *CatalogRepository) {
			err2 := cr.SetGoodQuantity("test-warehouse-ID", "test-ID", 7)
			_, presence := cr.GetWarehouses()["test-warehouse-ID"]
			assert.Equal(t, presence, true)
			assert.Equal(t, err2, nil)
		}),
	)
}
