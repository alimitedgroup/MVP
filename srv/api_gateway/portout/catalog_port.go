package portout

import (
	"github.com/alimitedgroup/MVP/common/dto"
)

type CatalogPortOut interface {
	// ListGoods recupera e ritorna la lista di merci registrate presso il servizio catalog
	ListGoods() (map[string]dto.Good, error)
	// ListStock recupera (se ho capito bene) per ogni merce
	// la somma delle disponibilità in ogni magazzino
	ListStock() (map[string]int64, error)
	// ListWarehouses ritorna la lista di magazzini esistenti
	ListWarehouses() (map[string]dto.Warehouse, error)
}
