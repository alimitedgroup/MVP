package business

import (
	"context"
	"log"

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
	transactionPort            port.ITransactionPort
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
	TransactionPort            port.ITransactionPort
}

func NewManageReservationService(p ManageReservationServiceParams) *ManageReservationService {
	return &ManageReservationService{
		p.CreateReservationEventPort, p.ApplyReservationEventPort, p.GetReservationPort,
		p.GetStockPort, p.CreateStockUpdatePort, p.IdempotentPort, p.Cfg, p.TransactionPort,
	}
}

func (s *ManageReservationService) CreateReservation(ctx context.Context, cmd port.CreateReservationCmd) (port.CreateReservationResponse, error) {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	validStock := true
	for _, good := range cmd.Goods {
		stock := s.getStockPort.GetFreeStock(model.GoodID(good.GoodID))
		log.Printf("free stock: good: %v stock: %v", good.GoodID, stock)
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

	err = s.applyReservationEventPort.ApplyReservationEvent(reservation)
	if err != nil {
		return port.CreateReservationResponse{}, err
	}

	idempotentCmd := port.IdempotentCmd{
		Event: "reservation",
		ID:    reservationId,
	}
	s.idempotentPort.SaveEventID(idempotentCmd)

	resp := port.CreateReservationResponse{
		ReservationID: reservationId,
	}
	return resp, nil
}

func (s *ManageReservationService) ApplyReservationEvent(cmd port.ApplyReservationEventCmd) error {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	goods := make([]model.ReservationGood, 0, len(cmd.Goods))

	for _, good := range cmd.Goods {
		goods = append(goods, model.ReservationGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	reserv := model.Reservation{
		ID:    cmd.ID,
		Goods: goods,
	}

	idempotentCmd := port.IdempotentCmd{
		Event: "reservation",
		ID:    cmd.ID,
	}
	if s.idempotentPort.IsAlreadyProcessed(idempotentCmd) {
		log.Printf("reservation already processed %v\n", cmd)
		return nil
	}

	err := s.applyReservationEventPort.ApplyReservationEvent(reserv)
	if err != nil {
		return err
	}

	return nil
}

func (s *ManageReservationService) ConfirmOrder(ctx context.Context, cmd port.ConfirmOrderCmd) error {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	log.Printf("ConfirmOrder: %v\n", cmd)
	if cmd.Status == "Filled" {
		for _, reserv := range cmd.Reservations {
			reservation, err := s.getReservationPort.GetReservation(model.ReservationID(reserv))
			if err != nil {
				continue
			}

			goods := make([]port.CreateStockUpdateGood, 0, len(reservation.Goods))
			for _, reservGood := range reservation.Goods {
				goodStock := s.getStockPort.GetStock(model.GoodID(reservGood.GoodID))

				goods = append(goods, port.CreateStockUpdateGood{
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

	} else if cmd.Status == "Cancelled" {
		for _, reserv := range cmd.Reservations {
			reservation, err := s.getReservationPort.GetReservation(model.ReservationID(reserv))
			if err != nil {
				continue
			}
			if err := s.applyReservationEventPort.ApplyOrderFilled(reservation); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ManageReservationService) ConfirmTransfer(ctx context.Context, cmd port.ConfirmTransferCmd) error {
	s.transactionPort.Lock()
	defer s.transactionPort.Unlock()

	if cmd.Status == "Filled" {
		if cmd.SenderID == s.cfg.ID {
			reservation, err := s.getReservationPort.GetReservation(model.ReservationID(cmd.ReservationID))
			if err != nil {
				// TODO: handle other errors
				return nil
			}

			goods := make([]port.CreateStockUpdateGood, 0, len(reservation.Goods))
			for _, reservGood := range reservation.Goods {
				goodStock := s.getStockPort.GetStock(model.GoodID(reservGood.GoodID))

				goods = append(goods, port.CreateStockUpdateGood{
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
			goods := make([]port.CreateStockUpdateGood, 0, len(cmd.Goods))

			for _, toAdd := range cmd.Goods {
				goodStock := s.getStockPort.GetStock(model.GoodID(toAdd.GoodID))

				goods = append(goods, port.CreateStockUpdateGood{
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
				ReservationID: cmd.ReservationID,
			}
			err := s.createStockUpdatePort.CreateStockUpdate(ctx, createCmd)
			if err != nil {
				return err
			}
		}
	} else if cmd.Status == "Cancelled" {
		reservation, err := s.getReservationPort.GetReservation(model.ReservationID(cmd.ReservationID))
		if err != nil {
			// TODO: handle other errors
			return nil
		}
		if err := s.applyReservationEventPort.ApplyOrderFilled(reservation); err != nil {
			return err
		}
	}

	return nil
}
