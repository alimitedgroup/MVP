package port

import model "github.com/alimitedgroup/MVP/srv/warehouse/model"

type UpdateStockUseCase interface {
	UpdateStock(UpdateStockCommand) error
}

type UpdateStockCommand struct {
	model.Good
}
