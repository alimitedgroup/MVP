package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/google/uuid"
	"go.uber.org/fx"
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

type ManageReservationServiceParams struct {
	fx.In

	CreateReservationEventPort port.IStoreReservationEventPort
	ApplyReservationEventPort  port.IApplyReservationEventPort
	GetReservationPort         port.IGetReservationPort
	GetStockPort               port.IGetStockPort
	CreateStockUpdatePort      port.ICreateStockUpdatePort
	IdempotentPort             port.IIdempotentPort
	Cfg                        *config.WarehouseConfig
}

func NewManageReservationService(p ManageReservationServiceParams) *ManageReservationService {
	return &ManageReservationService{
		p.CreateReservationEventPort, p.ApplyReservationEventPort, p.GetReservationPort,
		p.GetStockPort, p.CreateStockUpdatePort, p.IdempotentPort, p.Cfg,
	}
}

func (s *ManageReservationService) CreateReservation(ctx context.Context, cmd port.CreateReservationCmd) (port.CreateReservationResponse, error) {
	validStock := true
	for _, good := range cmd.Goods {
		stock := s.getStockPort.GetFreeStock(model.GoodID(good.GoodID))
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
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}
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
	goods := make([]model.ReservationGood, 0, len(cmd.Goods))

	for _, good := range cmd.Goods {
		goods = append(goods, model.ReservationGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	reserv := model.Reservation{
		ID:    cmd.Id,
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
			goodStock := s.getStockPort.GetStock(model.GoodID(reservGood.GoodID))

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
			TransferID:    "",
			ReservationID: reserv,
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

func (s *ManageReservationService) ConfirmTransfer(ctx context.Context, cmd port.ConfirmTransferCmd) error {
	if cmd.Status != "Filled" {
		return nil
	}

	if cmd.SenderID == s.cfg.ID {
		reservation, err := s.getReservationPort.GetReservation(model.ReservationId(cmd.ReservationId))
		if err != nil {
			// TODO: handle other errors
			return nil
		}

		goods := make([]port.CreateStockUpdateCmdGood, 0, len(reservation.Goods))
		for _, reservGood := range reservation.Goods {
			goodStock := s.getStockPort.GetStock(model.GoodID(reservGood.GoodID))

			goods = append(goods, port.CreateStockUpdateCmdGood{
				Good: model.GoodStock{
					ID:       reservGood.GoodID,
					Quantity: goodStock.Quantity - reservGood.Quantity,
				},
				QuantityDiff: reservGood.Quantity,
			})
		}

		createCmd := port.CreateStockUpdateCmd{
			Type:          port.CreateStockUpdateCmdTypeTransfer,
			Goods:         goods,
			OrderID:       "",
			TransferID:    cmd.TransferID,
			ReservationID: reservation.ID,
		}
		err = s.createStockUpdatePort.CreateStockUpdate(ctx, createCmd)
		if err != nil {
			return err
		}

		if err := s.applyReservationEventPort.ApplyOrderFilled(reservation); err != nil {
			return err
		}
	} else if cmd.ReceiverID == s.cfg.ID {
		goods := make([]port.CreateStockUpdateCmdGood, 0, len(cmd.Goods))

		for _, toAdd := range cmd.Goods {
			goodStock := s.getStockPort.GetStock(model.GoodID(toAdd.GoodID))

			goods = append(goods, port.CreateStockUpdateCmdGood{
				Good: model.GoodStock{
					ID:       toAdd.GoodID,
					Quantity: goodStock.Quantity + toAdd.Quantity,
				},
				QuantityDiff: toAdd.Quantity,
			})
		}

		createCmd := port.CreateStockUpdateCmd{
			Type:          port.CreateStockUpdateCmdTypeTransfer,
			Goods:         goods,
			OrderID:       "",
			TransferID:    cmd.TransferID,
			ReservationID: cmd.ReservationId,
		}
		err := s.createStockUpdatePort.CreateStockUpdate(ctx, createCmd)
		if err != nil {
			return err
		}
	}

	return nil
}
