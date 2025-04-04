package persistence

import (
	"sync"

	"github.com/alimitedgroup/MVP/common/dto"
)

type CatalogRepository struct {
	warehouseMap map[string]*dto.Warehouse
	goodStockMap map[string]int64
	mutex        sync.Mutex
}

func NewCatalogRepository() *CatalogRepository {
	return &CatalogRepository{warehouseMap: make(map[string]*dto.Warehouse), goodStockMap: make(map[string]int64)}
}

func (cr *CatalogRepository) GetGoodsGlobalQuantity() map[string]int64 {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	return cr.goodStockMap
}

func (cr *CatalogRepository) GetWarehouses() map[string]dto.Warehouse {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	result := make(map[string]dto.Warehouse)
	for key := range cr.warehouseMap {
		result[key] = *cr.warehouseMap[key]
	}
	return result
}

func (cr *CatalogRepository) SetGoodQuantity(warehouseID string, goodID string, newQuantity int64) error {
	/*
		Imposta la quantità di un bene in un magazzino e memorizza il nuovo stato globale della merce.
		Se il magazzino non esiste viene creato, se la merce non esiste viene memorizzata la quantità, ma non le info sulla merce)
	*/
	cr.addWarehouse(warehouseID)
	cr.mutex.Lock()
	_, presence := cr.goodStockMap[goodID]
	if !presence {
		//return catalogCommon.NewCustomError("Not a valid goodID")
		cr.goodStockMap[goodID] = newQuantity
		cr.warehouseMap[warehouseID].SetStock(goodID, newQuantity)
	} else {
		oldValue := cr.warehouseMap[warehouseID].GetGoodStock(goodID)
		delta := newQuantity - oldValue
		cr.warehouseMap[warehouseID].SetStock(goodID, newQuantity)
		cr.goodStockMap[goodID] += delta
	}
	cr.mutex.Unlock()
	return nil
}

func (cr *CatalogRepository) addWarehouse(warehouseID string) {
	/*
		Aggiunge un Warehouse alla lista dei Warehouse. Funzione invocata automaticamente quando l'aggiunta
		di una quantità di una merce determina l'assenza di un magazzino
	*/
	cr.mutex.Lock()
	_, presence := cr.warehouseMap[warehouseID]
	if presence {
		cr.mutex.Unlock()
		return
	}
	cr.warehouseMap[warehouseID] = dto.NewWarehouse(warehouseID)
	cr.mutex.Unlock()
}
