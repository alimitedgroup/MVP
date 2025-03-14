package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/google/uuid"
)

type PublishStockUpdateAdapter struct {
	broker       *broker.NatsMessageBroker
	warehouseCfg *config.WarehouseConfig
}

func NewPublishStockUpdateAdapter(broker *broker.NatsMessageBroker, warehouseCfg *config.WarehouseConfig) *PublishStockUpdateAdapter {
	return &PublishStockUpdateAdapter{broker, warehouseCfg}
}

func (a *PublishStockUpdateAdapter) CreateStockUpdate(ctx context.Context, cmd port.CreateStockUpdateCmd) error {
	stockUpdateId := uuid.New().String()

	goodsMsg := make([]stream.StockUpdateGood, 0, len(cmd.Goods))
	for _, v := range cmd.Goods {
		goodsMsg = append(goodsMsg, stream.StockUpdateGood{
			GoodID:   string(v.Good.ID),
			Quantity: v.Good.Quantity,
			Delta:    v.QuantityDiff,
		})
	}

	var stockUpdateType stream.StockUpdateType
	switch cmd.Type {
	case port.CreateStockUpdateCmdTypeAdd:
		stockUpdateType = stream.StockUpdateTypeAdd
	case port.CreateStockUpdateCmdTypeRemove:
		stockUpdateType = stream.StockUpdateTypeRemove
	default:
		return fmt.Errorf("unknown stock update type %s", cmd.Type)
	}

	streamMsg := stream.StockUpdate{
		ID:          stockUpdateId,
		WarehouseID: a.warehouseCfg.ID,
		Goods:       goodsMsg,
		TransferID:  "",
		OrderID:     "",
		Type:        stockUpdateType,
	}

	payload, err := json.Marshal(streamMsg)
	if err != nil {
		return err
	}

	resp, err := a.broker.Js.Publish(ctx, fmt.Sprintf("stock.update.%s", a.warehouseCfg.ID), payload)
	if err != nil {
		return err
	}

	_ = resp

	return nil
}
