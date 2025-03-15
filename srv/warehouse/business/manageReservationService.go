package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/google/uuid"
)

type ManageReservationService struct {
	createReservationEventPort port.IStoreReservationEventPort
	applyReservationEventPort  port.IApplyReservationEventPort
}

func NewManageReservationService(
	createReservationEventPort port.IStoreReservationEventPort,
	applyReservationEventPort port.IApplyReservationEventPort,
) *ManageReservationService {
	return &ManageReservationService{createReservationEventPort, applyReservationEventPort}
}

func (s *ManageReservationService) CreateReservation(ctx context.Context, cmd port.CreateReservationCmd) (port.CreateReservationResponse, error) {
	goods := make([]model.ReservationGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.ReservationGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}
	reservationId := uuid.New().String()
	reservation := model.Reservation{ID: reservationId, Goods: goods}

	err := s.createReservationEventPort.StoreReservationEvent(ctx, reservation)
	if err != nil {
		return port.CreateReservationResponse{}, err
	}

	resp := port.CreateReservationResponse{
		ReservationID: reservationId,
	}
	return resp, nil
}

func (s *ManageReservationService) ApplyReservationEvent(cmd port.ApplyReservationEventCmd) error {
	return nil
}
