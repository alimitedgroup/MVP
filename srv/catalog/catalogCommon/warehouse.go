package catalogCommon

type Warehouse struct {
	ID    string
	stock map[string]int64
}

func NewWarehouse(ID string) *Warehouse {
	return &Warehouse{ID, make(map[string]int64)}
}

func (w *Warehouse) SetStock(ID string, newQuantity int64) {
	_, presence := w.stock[ID]
	if !presence {
		w.addGood(ID)
	}
	w.stock[ID] = newQuantity
}

func (w *Warehouse) addGood(ID string) {
	w.stock[ID] = 0
}

func (w *Warehouse) GetGoodStock(id string) int64 {
	value, presence := w.stock[id]
	if !presence {
		return int64(0)
	}
	return value
}
