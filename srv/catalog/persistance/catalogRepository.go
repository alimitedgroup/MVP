package persistance

import "github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"

type catalogRepository struct {
	warehouseMap map[string]*warehouse
	goodMap      map[string]*catalogCommon.Good
	goodStockMap map[string]int64
}

func NewCatalogRepository() *catalogRepository {
	return &catalogRepository{warehouseMap: make(map[string]*warehouse), goodMap: make(map[string]*catalogCommon.Good), goodStockMap: make(map[string]int64)}
}

func (cr *catalogRepository) GetGoods() map[string]catalogCommon.Good {
	result := make(map[string]catalogCommon.Good)
	for key := range cr.goodMap {
		result[key] = *cr.goodMap[key]
	}
	return result
}

func (cr *catalogRepository) GetGoodsGlobalQuantity() map[string]int64 {
	return cr.goodStockMap
}

func (cr *catalogRepository) GetWarehouses() map[string]warehouse {
	result := make(map[string]warehouse)
	for key := range cr.warehouseMap {
		result[key] = *cr.warehouseMap[key]
	}
	return result
}

func (cr *catalogRepository) SetGoodQuantity(warehouseID string, goodID string, newQuantity int64) error {
	/*
		Imposta la quantità di un bene in un magazzino e memorizza il nuovo stato globale della merce.
		Se il magazzino non esiste viene creato, se la merce non esiste ritorna un errore (le informazioni
		attuali non bastano per creare autonomamente la nuova merce)
	*/
	cr.addWarehouse(warehouseID)
	_, presence := cr.goodMap[goodID]
	if !presence {
		return catalogCommon.NewCustomError("Not a valid goodID")
	}
	oldValue := cr.warehouseMap[warehouseID].GetGoodStock(goodID)
	delta := newQuantity - oldValue
	cr.warehouseMap[warehouseID].SetStock(goodID, newQuantity)
	cr.goodStockMap[goodID] += delta
	return nil
}

func (cr *catalogRepository) addWarehouse(warehouseID string) {
	/*
		Aggiunge un Warehouse alla lista dei Warehouse. Funzione invocata automaticamente quando l'aggiunta
		di una quantità di una merce determina l'assenza di un magazzino
	*/
	_, presence := cr.warehouseMap[warehouseID]
	if presence {
		return
	}
	cr.warehouseMap[warehouseID] = NewWarehouse(warehouseID)
}

func (cr *catalogRepository) AddGood(goodID string, name string, description string) error {
	_, presence := cr.goodMap[goodID]
	if presence {
		return cr.changeGoodData(goodID, name, description)
	}
	cr.goodMap[goodID] = catalogCommon.NewGood(goodID, name, description)
	return nil
}

func (cr *catalogRepository) changeGoodData(goodID string, newName string, newDescription string) error {
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
