package controller

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/alimitedgroup/MVP/common/stream"
	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/notification/service/portin"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

type notificationController struct {
	addStockUpdateUseCase serviceportin.IAddStockUpdateUseCase
	addQueryRuleUseCase   serviceportin.IAddQueryRuleUseCase
}

type NotificationControllerParams struct {
	fx.In
	AddStockUpdateUseCase serviceportin.IAddStockUpdateUseCase
	AddQueryRuleUseCase   serviceportin.IAddQueryRuleUseCase
}

func NewNotificationController(p NotificationControllerParams) *notificationController {
	return &notificationController{
		addStockUpdateUseCase: p.AddStockUpdateUseCase,
		addQueryRuleUseCase:   p.AddQueryRuleUseCase,
	}
}

func (nc *notificationController) addStockUpdateRequest(ctx context.Context, msg jetstream.Msg) error {
	request := &stream.StockUpdate{}

	err := json.Unmarshal(msg.Data(), request)
	if err != nil {
		return err
	}

	goods := make([]servicecmd.StockGood, len(request.Goods))
	for i, g := range request.Goods {
		goods[i] = servicecmd.StockGood{
			ID:       g.GoodID,
			Quantity: int(g.Quantity),
			Delta:    int(g.Delta),
		}
	}

	cmd := servicecmd.NewAddStockUpdateCmd(request.WarehouseID, string(request.Type), request.OrderID, request.TransferID, goods, time.Now().Unix())
	_, err = nc.addStockUpdateUseCase.AddStockUpdate(cmd)

	return err
}


func (nc *notificationController) addQueryRuleRequest(ctx context.Context, msg jetstream.Msg) error {
	log.Printf("addQueryRuleRequest ricevuto: %s", string(msg.Data()))

	request := &stream.AddQueryRule{}

	err := json.Unmarshal(msg.Data(), request)
	if err != nil {
		return err
	}

	cmd := servicecmd.NewAddQueryRuleCmd(request.GoodID, request.Operator, request.Threshold)
	response := nc.addQueryRuleUseCase.AddQueryRule(cmd)

	return response.GetOperationResult()
}
