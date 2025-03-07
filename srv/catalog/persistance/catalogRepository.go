package persistance

import (
	"sync"

	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
)

type CatalogRepository struct {
	warehouseMap map[string]*catalogCommon.Warehouse
	goodMap      map[string]*catalogCommon.Good
	goodStockMap map[string]int64
}

func NewCatalogRepository() *CatalogRepository {
	return &CatalogRepository{warehouseMap: make(map[string]*catalogCommon.Warehouse), goodMap: make(map[string]*catalogCommon.Good), goodStockMap: make(map[string]int64)}
}

func (cr *CatalogRepository) GetGoods() map[string]catalogCommon.Good {
	result := make(map[string]catalogCommon.Good)
	for key := range cr.goodMap {
		result[key] = *cr.goodMap[key]
	}
	return result
}

func (cr *CatalogRepository) GetGoodsGlobalQuantity() map[string]int64 {
	return cr.goodStockMap
}

func (cr *CatalogRepository) GetWarehouses() map[string]catalogCommon.Warehouse {
	result := make(map[string]catalogCommon.Warehouse)
	for key := range cr.warehouseMap {
		result[key] = *cr.warehouseMap[key]
	}
	return result
}

var mutex sync.Mutex

func (cr *CatalogRepository) SetGoodQuantity(warehouseID string, goodID string, newQuantity int64) error {
	/*
		Imposta la quantità di un bene in un magazzino e memorizza il nuovo stato globale della merce.
		Se il magazzino non esiste viene creato, se la merce non esiste viene memorizzata la quantità, ma non le info sulla merce)
	*/
	mutex.Lock()
	cr.addWarehouse(warehouseID)
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
	mutex.Unlock()
	return nil
}

func (cr *CatalogRepository) addWarehouse(warehouseID string) {
	/*
		Aggiunge un Warehouse alla lista dei Warehouse. Funzione invocata automaticamente quando l'aggiunta
		di una quantità di una merce determina l'assenza di un magazzino
	*/
	_, presence := cr.warehouseMap[warehouseID]
	if presence {
		return
	}
	cr.warehouseMap[warehouseID] = catalogCommon.NewWarehouse(warehouseID)
}

func (cr *CatalogRepository) AddGood(goodID string, name string, description string) error {
	_, presence := cr.goodMap[goodID]
	if presence {
		return cr.changeGoodData(goodID, name, description)
	}
	cr.goodMap[goodID] = catalogCommon.NewGood(goodID, name, description)
	return nil
}

func (cr *CatalogRepository) changeGoodData(goodID string, newName string, newDescription string) error {
	_, presence := cr.goodMap[goodID]
	if !presence {
		return catalogCommon.NewCustomError("Not a valid goodID")
	}
	err := cr.goodMap[goodID].SetName(newName)
	if err != nil {
		return err
	}
	err = cr.goodMap[goodID].SetDescription(newDescription)
	if err != nil {
		return err
	}
	return nil
}
