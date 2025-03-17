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
	idempotentPort             port.IIdempotentPort
	cfg                        *config.WarehouseConfig
}

func NewManageReservationService(
	createReservationEventPort port.IStoreReservationEventPort,
	applyReservationEventPort port.IApplyReservationEventPort,
	getReservationPort port.IGetReservationPort,
	getStockPort port.IGetStockPort,
	createStockUpdatePort port.ICreateStockUpdatePort,
	idempotentPort port.IIdempotentPort,
	cfg *config.WarehouseConfig,
) *ManageReservationService {
	return &ManageReservationService{
		createReservationEventPort, applyReservationEventPort, getReservationPort,
		getStockPort, createStockUpdatePort, idempotentPort, cfg,
	}
}

func (s *ManageReservationService) CreateReservation(ctx context.Context, cmd port.CreateReservationCmd) (port.CreateReservationResponse, error) {
	validStock := true
	for _, good := range cmd.Goods {
		stock := s.getStockPort.GetFreeStock(model.GoodId(good.GoodID))
		if stock.Quantity < good.Quantity {
			validStock = false
			break
		}
	}

	if !validStock {
		return port.CreateReservationResponse{}, port.ErrNotEnoughStock
	}

	reservationId := uuid.New().String()
	goods := make([]model.ReservationGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, model.ReservationGood{
			GoodID:   model.GoodId(good.GoodID),
			Quantity: good.Quantity,
		})
	}
	reservation := model.Reservation{ID: model.ReservationId(reservationId), Goods: goods}
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

	idempotentCmd := port.IdempotentCmd{
		Event: "reservation",
		Id:    cmd.Id,
	}
	if s.idempotentPort.IsAlreadyProcessed(idempotentCmd) {
		return nil
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

		if err := s.applyReservationEventPort.ApplyOrderFilled(reservation); err != nil {
			return err
		}
	}

	return nil
}
