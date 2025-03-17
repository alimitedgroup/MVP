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
			GoodID:   good.GoodId,
			Quantity: good.Quantity,
		})
	}
	streamMsg := stream.OrderUpdate{
		ID:           cmd.ID,
		Status:       cmd.Status,
		Name:         cmd.Name,
		Email:        cmd.Email,
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
			ID:       model.GoodId(good.GoodId),
			Quantity: good.Quantity,
		})
	}

	return model.Order{
		Id:           model.OrderID(cmd.ID),
		Status:       cmd.Status,
		Name:         cmd.Name,
		Email:        cmd.Email,
		Address:      cmd.Address,
		Goods:        modelGoods,
		Reservations: cmd.Reservations,
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

	goods := make([]internalStream.ContactWarehousesGood, 0, len(cmd.Order.Goods))
	for _, good := range cmd.Order.Goods {
		goods = append(goods, internalStream.ContactWarehousesGood{
			GoodId:   string(good.ID),
			Quantity: good.Quantity,
		})
	}

	streamMsg := internalStream.ContactWarehouses{
		Order: internalStream.ContactWarehousesOrder{
			ID:           string(cmd.Order.Id),
			Status:       cmd.Order.Status,
			Name:         cmd.Order.Name,
			Email:        cmd.Order.Email,
			Address:      cmd.Order.Address,
			UpdateTime:   cmd.Order.UpdateTime,
			CreationTime: cmd.Order.CreationTime,
			Goods:        goods,
			Reservations: cmd.Order.Reservations,
		},
		TransferId:            cmd.TransferId,
		LastContact:           cmd.LastContact,
		ConfirmedReservations: confirmed,
		ExcludeWarehouses:     cmd.ExcludeWarehouses,
	}

	payload, err := json.Marshal(streamMsg)
	if err != nil {
		return err
	}

	resp, err := a.broker.Js.Publish(ctx, "order.contact.warehouses", payload)
	if err != nil {
		return err
	}

	_ = resp

	return nil
}

func (a *NatsStreamAdapter) RequestReservation(ctx context.Context, cmd port.RequestReservationCmd) (port.RequestReservationResponse, error) {
	goods := make([]request.ReserveStockItem, 0, len(cmd.Items))
	for _, good := range cmd.Items {
		goods = append(goods, request.ReserveStockItem{
			GoodID:   good.GoodId,
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

	return port.RequestReservationResponse{Id: respDto.Message.ReservationID}, nil
}
