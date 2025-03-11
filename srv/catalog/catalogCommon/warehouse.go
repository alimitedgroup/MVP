package catalogCommon

type Warehouse struct {
	ID    string           `json:"id"`
	Stock map[string]int64 `json:"stock"`
}

func NewWarehouse(ID string) *Warehouse {
	return &Warehouse{ID, make(map[string]int64)}
}

func (w *Warehouse) SetStock(ID string, newQuantity int64) {
	_, presence := w.Stock[ID]
	if !presence {
		w.addGood(ID)
	}
	w.Stock[ID] = newQuantity
}

func (w *Warehouse) addGood(ID string) {
	w.Stock[ID] = 0
}

func (w *Warehouse) GetGoodStock(id string) int64 {
	value, presence := w.Stock[id]
	if !presence {
		return int64(0)
	}
	return value
}
