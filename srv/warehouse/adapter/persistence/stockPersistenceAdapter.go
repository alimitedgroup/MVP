package persistence

import (
	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
)

type StockPersistanceAdapter struct {
	stockRepo IStockRepository
}

func NewStockPersistanceAdapter(stockRepo IStockRepository) *StockPersistanceAdapter {
	return &StockPersistanceAdapter{stockRepo}
}

func (s *StockPersistanceAdapter) ApplyStockUpdate(goods []model.GoodStock) {
	for _, good := range goods {
		s.stockRepo.SetStock(string(good.ID), good.Quantity)
	}
}

func (s *StockPersistanceAdapter) ApplyReservationEvent(reservation model.Reservation) error {
	for _, good := range reservation.Goods {
		if err := s.stockRepo.ReserveStock(string(reservation.ID), string(good.GoodID), good.Quantity); err != nil {
			return err
		}
	}

	return nil
}

func (s *StockPersistanceAdapter) ApplyOrderFilled(reservation model.Reservation) error {
	for _, good := range reservation.Goods {
		if err := s.stockRepo.UnReserveStock(string(good.GoodID), good.Quantity); err != nil {
			return err
		}
	}

	return nil
}

func (s *StockPersistanceAdapter) GetStock(goodId model.GoodID) model.GoodStock {
	stock := s.stockRepo.GetStock(string(goodId))
	return model.GoodStock{
		ID:       goodId,
		Quantity: stock,
	}
}

func (s *StockPersistanceAdapter) GetFreeStock(goodId model.GoodID) model.GoodStock {
	stock := s.stockRepo.GetFreeStock(string(goodId))
	return model.GoodStock{
		ID:       goodId,
		Quantity: stock,
	}
}

func (s *StockPersistanceAdapter) GetReservation(reservationId model.ReservationId) (model.Reservation, error) {
	reserv, err := s.stockRepo.GetReservation(string(reservationId))
	if err != nil {
		return model.Reservation{}, err
	}

	goods := make([]model.ReservationGood, 0, len(reserv.Goods))

	for goodId, qty := range reserv.Goods {
		goods = append(goods, model.ReservationGood{
			GoodID:   model.GoodID(goodId),
			Quantity: qty,
		})
	}

	return model.Reservation{
		ID:    reservationId,
		Goods: goods,
	}, nil

}
