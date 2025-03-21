package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	internalStream "github.com/alimitedgroup/MVP/srv/order/adapter/stream"
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type NatsStreamAdapter struct {
	broker *broker.NatsMessageBroker
}

func NewNatsStreamAdapter(broker *broker.NatsMessageBroker) *NatsStreamAdapter {
	return &NatsStreamAdapter{broker}
}

func (a *NatsStreamAdapter) SendOrderUpdate(ctx context.Context, cmd port.SendOrderUpdateCmd) (model.Order, error) {
	now := time.Now()

	updateTime := now.UnixMilli()
	var creationTime int64
	if cmd.CreationTime == 0 {
		creationTime = updateTime
	} else {
		creationTime = cmd.CreationTime
	}

	goods := make([]stream.OrderUpdateGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, stream.OrderUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}
	streamMsg := stream.OrderUpdate{
		ID:           cmd.ID,
		Status:       cmd.Status,
		Name:         cmd.Name,
		FullName:     cmd.FullName,
		Address:      cmd.Address,
		Goods:        goods,
		Reservations: cmd.Reservations,
		CreationTime: creationTime,
		UpdateTime:   updateTime,
	}

	payload, err := json.Marshal(streamMsg)
	if err != nil {
		return model.Order{}, err
	}

	resp, err := a.broker.Js.Publish(ctx, "order.update", payload)
	if err != nil {
		return model.Order{}, err
	}

	_ = resp

	modelGoods := make([]model.GoodStock, 0, len(goods))
	for _, good := range cmd.Goods {
		modelGoods = append(modelGoods, model.GoodStock{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	return model.Order{
		ID:           cmd.ID,
		Status:       cmd.Status,
		Name:         cmd.Name,
		FullName:     cmd.FullName,
		Address:      cmd.Address,
		Goods:        modelGoods,
		Reservations: cmd.Reservations,
		CreationTime: creationTime,
		UpdateTime:   updateTime,
	}, nil
}

func (a *NatsStreamAdapter) SendTransferUpdate(ctx context.Context, cmd port.SendTransferUpdateCmd) (model.Transfer, error) {
	now := time.Now()

	updateTime := now.UnixMilli()
	var creationTime int64
	if cmd.CreationTime == 0 {
		creationTime = updateTime
	} else {
		creationTime = cmd.CreationTime
	}

	goods := make([]stream.TransferUpdateGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, stream.TransferUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}
	streamMsg := stream.TransferUpdate{
		ID:            cmd.ID,
		Status:        cmd.Status,
		SenderID:      cmd.SenderID,
		ReceiverID:    cmd.ReceiverID,
		ReservationId: cmd.ReservationId,
		Goods:         goods,
		CreationTime:  creationTime,
		UpdateTime:    updateTime,
	}

	payload, err := json.Marshal(streamMsg)
	if err != nil {
		return model.Transfer{}, err
	}

	resp, err := a.broker.Js.Publish(ctx, "transfer.update", payload)
	if err != nil {
		return model.Transfer{}, err
	}

	_ = resp

	modelGoods := make([]model.GoodStock, 0, len(goods))
	for _, good := range cmd.Goods {
		modelGoods = append(modelGoods, model.GoodStock{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	return model.Transfer{
		ID:           cmd.ID,
		Status:       cmd.Status,
		SenderId:     cmd.SenderID,
		ReceiverId:   cmd.ReceiverID,
		Goods:        modelGoods,
		CreationTime: creationTime,
		UpdateTime:   updateTime,
	}, nil
}

func (a *NatsStreamAdapter) SendContactWarehouses(ctx context.Context, cmd port.SendContactWarehouseCmd) error {
	confirmed := make([]internalStream.ConfirmedReservation, 0, len(cmd.ConfirmedReservations))

	for _, reservation := range cmd.ConfirmedReservations {
		confirmed = append(confirmed, internalStream.ConfirmedReservation{
			WarehouseId:   reservation.WarehouseId,
			ReservationID: reservation.ReservationID,
			Goods:         reservation.Goods,
		})
	}

	var transfer *internalStream.ContactWarehousesTransfer
	var order *internalStream.ContactWarehousesOrder
	if cmd.Type == port.SendContactWarehouseTypeOrder {
		goods := make([]internalStream.ContactWarehousesGood, 0, len(cmd.Order.Goods))
		for _, good := range cmd.Order.Goods {
			goods = append(goods, internalStream.ContactWarehousesGood{
				GoodID:   string(good.GoodID),
				Quantity: good.Quantity,
			})
		}
		order = &internalStream.ContactWarehousesOrder{
			ID:           string(cmd.Order.ID),
			Status:       cmd.Order.Status,
			Name:         cmd.Order.Name,
			FullName:     cmd.Order.FullName,
			Address:      cmd.Order.Address,
			UpdateTime:   cmd.Order.UpdateTime,
			CreationTime: cmd.Order.CreationTime,
			Goods:        goods,
			Reservations: cmd.Order.Reservations,
		}

	} else if cmd.Type == port.SendContactWarehouseTypeTransfer {
		goods := make([]internalStream.ContactWarehousesGood, 0, len(cmd.Transfer.Goods))
		for _, good := range cmd.Transfer.Goods {
			goods = append(goods, internalStream.ContactWarehousesGood{
				GoodID:   string(good.GoodID),
				Quantity: good.Quantity,
			})
		}
		transfer = &internalStream.ContactWarehousesTransfer{
			ID:            string(cmd.Transfer.ID),
			Status:        cmd.Transfer.Status,
			SenderId:      string(cmd.Transfer.SenderId),
			ReceiverId:    string(cmd.Transfer.ReceiverId),
			UpdateTime:    cmd.Transfer.UpdateTime,
			CreationTime:  cmd.Transfer.CreationTime,
			Goods:         goods,
			ReservationId: cmd.Transfer.ReservationID,
		}
	}

	streamMsg := internalStream.ContactWarehouses{
		Type:                  internalStream.ContactWarehousesType(cmd.Type),
		Order:                 order,
		Transfer:              transfer,
		ConfirmedReservations: confirmed,
		ExcludeWarehouses:     cmd.ExcludeWarehouses,
		RetryInTime:           cmd.RetryInTime,
		RetryUntil:            cmd.RetryUntil,
	}

	payload, err := json.Marshal(streamMsg)
	if err != nil {
		return err
	}

	resp, err := a.broker.Js.Publish(ctx, "contact.warehouses", payload)
	if err != nil {
		return err
	}

	_ = resp

	return nil
}

func (a *NatsStreamAdapter) RequestReservation(ctx context.Context, cmd port.RequestReservationCmd) (port.RequestReservationResponse, error) {
	goods := make([]request.ReserveStockItem, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, request.ReserveStockItem{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	reqMsg := request.ReserveStockRequestDTO{
		Goods: goods,
	}

	payload, err := json.Marshal(reqMsg)
	if err != nil {
		return port.RequestReservationResponse{}, err
	}

	resp, err := a.broker.Nats.Request(fmt.Sprintf("warehouse.%s.reservation.create", cmd.WarehouseId), payload, 5*time.Second)
	if err != nil {
		return port.RequestReservationResponse{}, err
	}

	var respDto response.ReserveStockResponseDTO
	if err := json.Unmarshal(resp.Data, &respDto); err != nil {
		return port.RequestReservationResponse{}, err
	}

	if respDto.Error != "" {
		return port.RequestReservationResponse{}, port.ErrNotEnoughStock
	}

	return port.RequestReservationResponse{Id: respDto.Message.ReservationID}, nil
}
