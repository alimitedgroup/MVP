package goodRepository

type catalogRepository struct {
	warehouseMap map[string]*warehouse
	goodMap      map[string]*good
	goodStockMap map[string]int64
}

func NewCatalogRepository() *catalogRepository {
	return &catalogRepository{warehouseMap: make(map[string]*warehouse), goodMap: make(map[string]*good)}
}

func (cr *catalogRepository) GetGoods() map[string]good {
	result := make(map[string]good)
	for key := range cr.goodMap {
		result[key] = *cr.goodMap[key]
	}
	return result
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
		return CustomError{"Not a valid good ID"}
	}
	cr.warehouseMap[warehouseID].SetStock(goodID, newQuantity)
	cr.goodMap[goodID].SetQuantity(newQuantity)
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

func (cr *catalogRepository) AddGood(name string, description string, goodID string) error {
	_, presence := cr.goodMap[goodID]
	if presence {
		return CustomError{"Provided goodID already exists"}
	}
	cr.goodMap[goodID] = NewGood(goodID, name, description, 0)
	return nil
}
