package goodRepository

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

func TestNotAValidGoodID_ChangeData(t *testing.T) {
	// Se si modificano le informazioni di una merce non esistente ritorna un errore
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *catalogRepository) {
			cr.AddGood("test-name", "test-desc", "test-ID")
			err := cr.ChangeGoodData("ciao", "test-name", "test-desc2")
			assert.Equal(t, err.Error(), "Not a valid goodID")
		}),
	)
}

func TestNotAValidGoodID_SetQt(t *testing.T) {
	// Se si modifica la quantità di una merce non esistente ritorna un errore
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *catalogRepository) {
			cr.AddGood("test-name", "test-desc", "test-ID")
			err := cr.SetGoodQuantity("test-name2", "un bell'ID", 7)
			assert.Equal(t, err.Error(), "Not a valid goodID")
		}),
	)
}
func TestAlreadyExistentGood(t *testing.T) {
	// Se si prova ad aggiungere una merce con ID già esistente restituisce un errore
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *catalogRepository) {
			cr.AddGood("test-name", "test-desc", "test-ID")
			err := cr.AddGood("test-name2", "test-desc2", "test-ID")
			assert.Equal(t, err.Error(), "Provided goodID already exists")
		}),
	)
}

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
			assert.Equal(t, cr.GetGoodsGlobalQuantity()["test-ID"], int64(7))
		}),
	)
}

func TestChangeGoodData(t *testing.T) {
	// Testa la modifica di una merce
	fx.New(
		fx.Provide(NewCatalogRepository),
		fx.Invoke(func(cr *catalogRepository) {
			cr.AddGood("test-name", "test-description", "test-ID")
			cr.SetGoodQuantity("test-warehouse-ID", "test-ID", 7)
			cr.ChangeGoodData("test-ID", "new-test-name", "new-test-description")
			resultGoods := cr.GetGoods()

			assert.Equal(t, resultGoods["test-ID"].GetID(), "test-ID")
			assert.Equal(t, resultGoods["test-ID"].GetName(), "new-test-name")
			assert.Equal(t, resultGoods["test-ID"].GetDescription(), "new-test-description")
			assert.Equal(t, cr.GetGoodsGlobalQuantity()["test-ID"], int64(7))
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
