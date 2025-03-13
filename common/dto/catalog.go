package dto

import "errors"

type Good struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"` //it is a converted uuid.UUID
}

type Warehouse struct {
	ID    string           `json:"id"`
	Stock map[string]int64 `json:"stock"`
}

var (
	ErrEmptyDescription = errors.New("description is empty")
	ErrEmptyName        = errors.New("name is empty")
)

func NewGood(ID string, name string, description string) *Good {
	return &Good{name, description, ID}
}

func (g Good) GetID() string {
	return g.ID
}

func (g Good) GetName() string {
	return g.Name
}

func (g Good) GetDescription() string {
	return g.Description
}

func (g *Good) SetDescription(newDescription string) error {
	if newDescription == "" {
		return ErrEmptyDescription
	}
	g.Description = newDescription
	return nil
}

func (g *Good) SetName(newName string) error {
	if newName == "" {
		return ErrEmptyName
	}
	g.Name = newName
	return nil
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

func NewWarehouse(ID string) *Warehouse {
	return &Warehouse{ID, make(map[string]int64)}
}

type GetGoodsDataResponseDTO struct {
	GoodMap map[string]Good `json:"goods"`
	Err     string          `json:"error"`
}

type GetWarehouseResponseDTO struct {
	WarehouseMap map[string]Warehouse `json:"warehouse_map"`
	Err          string               `json:"error"`
}

type GetGoodsQuantityResponseDTO struct {
	GoodMap map[string]int64 `json:"goods"`
	Err     string           `json:"error"`
}
