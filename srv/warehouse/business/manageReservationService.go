package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/google/uuid"
)

type ManageReservationService struct {
	createReservationEventPort port.IStoreReservationEventPort
	applyReservationEventPort  port.IApplyReservationEventPort
	getReservationPort         port.IGetReservationPort
	getStockPort               port.IGetStockPort
	createStockUpdatePort      port.ICreateStockUpdatePort
	cfg                        *config.WarehouseConfig
}

func NewManageReservationService(
	createReservationEventPort port.IStoreReservationEventPort,
	applyReservationEventPort port.IApplyReservationEventPort,
	getReservationPort port.IGetReservationPort,
	getStockPort port.IGetStockPort,
	createStockUpdatePort port.ICreateStockUpdatePort,
	cfg *config.WarehouseConfig,
) *ManageReservationService {
	return &ManageReservationService{
		createReservationEventPort, applyReservationEventPort, getReservationPort,
		getStockPort, createStockUpdatePort, cfg,
	}
}

func (s *ManageReservationService) CreateReservation(ctx context.Context, cmd port.CreateReservationCmd) (port.CreateReservationResponse, error) {
	goods := make([]model.ReservationGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.ReservationGood{
			GoodID:   model.GoodId(good.GoodID),
			Quantity: good.Quantity,
		})
	}
	reservationId := uuid.New().String()
	reservation := model.Reservation{ID: model.ReservationId(reservationId), Goods: goods}

	// TODO: check if the goods are available

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
	goods := make([]model.ReservationGood, 0, len(cmd.Goods))

	for _, good := range cmd.Goods {
		goods = append(goods, model.ReservationGood{
			GoodID:   model.GoodId(good.GoodID),
			Quantity: good.Quantity,
		})
	}

	reserv := model.Reservation{
		ID:    model.ReservationId(cmd.Id),
		Goods: goods,
	}

	err := s.applyReservationEventPort.ApplyReservationEvent(reserv)
	if err != nil {
		return err
	}

	return nil
}

func (s *ManageReservationService) ConfirmOrder(ctx context.Context, cmd port.ConfirmOrderCmd) error {
	if cmd.Status != "Filled" {
		return nil
	}

	for _, reserv := range cmd.Reservations {
		reservation, err := s.getReservationPort.GetReservation(model.ReservationId(reserv))
		if err != nil {
			continue
		}

		goods := make([]port.CreateStockUpdateCmdGood, 0, len(reservation.Goods))
		for _, reservGood := range reservation.Goods {
			goodStock := s.getStockPort.GetStock(reservGood.GoodID)

			goods = append(goods, port.CreateStockUpdateCmdGood{
				Good: model.GoodStock{
					ID:       reservGood.GoodID,
					Quantity: goodStock.Quantity - reservGood.Quantity,
				},
				QuantityDiff: reservGood.Quantity,
			})
		}

		createCmd := port.CreateStockUpdateCmd{
			Type:          port.CreateStockUpdateCmdTypeOrder,
			Goods:         goods,
			OrderID:       cmd.OrderID,
			ReservationID: reserv,
			TransferID:    "",
		}
		err = s.createStockUpdatePort.CreateStockUpdate(ctx, createCmd)
		if err != nil {
			return err
		}

		// TODO: unreserve the goods
	}

	return nil
}
