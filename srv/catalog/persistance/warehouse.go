package goodRepository

type warehouse struct {
	ID    string
	stock map[string]int64
}

func NewWarehouse(ID string) *warehouse {
	return &warehouse{ID, make(map[string]int64)}
}

func (w *warehouse) SetStock(ID string, newQuantity int64) {
	_, presence := w.stock[ID]
	if !presence {
		w.addGood(ID)
	}
	w.stock[ID] = newQuantity
}

func (w *warehouse) addGood(ID string) {
	w.stock[ID] = 0
}

func (w *warehouse) GetStock() map[string]int64 {
	return w.stock
}
