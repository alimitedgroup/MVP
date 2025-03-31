package adapterout

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/nats-io/nats.go"
)

type NotificationAdapter struct {
	influxClient influxdb2.Client
	brk          *broker.NatsMessageBroker
	ruleRepo     portout.RuleRepository
	writeApi     api.WriteAPI
	queryApi     api.QueryAPI
	*zap.Logger
}

func NewNotificationAdapter(influxClient influxdb2.Client, brk *broker.NatsMessageBroker, ruleRepo portout.RuleRepository, logger *zap.Logger) *NotificationAdapter {
	return &NotificationAdapter{
		influxClient: influxClient,
		brk:          brk,
		ruleRepo:     ruleRepo,
		writeApi:     influxClient.WriteAPI("my-org", "stockdb"),
		queryApi:     influxClient.QueryAPI("my-org"),
		Logger:       logger,
	}
}

// =========== StockRepository port-out ===========

func (na *NotificationAdapter) SaveStockUpdate(cmd *types.AddStockUpdateCmd) error {
	respmsg, err := na.brk.Nats.Request("catalog.getGoodsGlobalQuantity", []byte("{}"), nats.DefaultTimeout)
	if err != nil {
		na.Error("Error querying catalog.getGoodsGlobalQuantity", zap.Error(err))
		return err
	}

	var resp dto.GetGoodsQuantityResponseDTO
	err = json.Unmarshal(respmsg.Data, &resp)
	if err != nil {
		na.Error("Error querying catalog.getGoodsGlobalQuantity", zap.Error(err))
		return err
	}

	for _, good := range cmd.Goods {
		na.writeApi.WritePoint(influxdb2.NewPoint(
			"stock_measurement",
			map[string]string{"warehouse_id": cmd.WarehouseID, "good_id": good.ID},
			map[string]interface{}{"quantity": int64(good.Delta) + resp.GoodMap[good.ID]},
			time.Now(),
		))
	}

	return nil
}

// =========== StockEventPublisher port-out ===========

func (na *NotificationAdapter) PublishStockAlert(alert types.StockAlertEvent) error {
	serviceId, exist := os.LookupEnv("ENV_SERVICE_ID")
	if !exist {
		serviceId = "DEFAULT"
	}
	data, err := json.Marshal(alert)
	if err != nil {
		na.Error("Error marshalling alert", zap.Error(err))
		return err
	}
	if err = na.brk.Nats.Publish(fmt.Sprintf("stock.alert.%s", serviceId), data); err != nil {
		na.Error("Error publishing alert to NATS", zap.Error(err))
		return err
	}
	return nil
}

// =========== RuleQueryRepository port-out ===========

func (na *NotificationAdapter) GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse {
	fluxQuery := `from(bucket:"stockdb")
		|> range(start:-7d)
		|> filter(fn:(r)=> r["_measurement"]=="stock_measurement")
		|> filter(fn:(r)=> r["good_id"]=="` + goodID + `")
		|> filter(fn:(r)=> r["_field"]=="quantity")
		|> last()`
	result, err := na.queryApi.Query(context.Background(), fluxQuery)
	if err != nil {
		return types.NewGetRuleResultResponse(goodID, 0, err)
	}
	for result.Next() {
		val, ok := result.Record().Value().(int64)
		if !ok {
			continue
		}
		return types.NewGetRuleResultResponse(goodID, int(val), nil)
	}
	if result.Err() != nil {
		return types.NewGetRuleResultResponse(goodID, 0, result.Err())
	}
	return types.NewGetRuleResultResponse(goodID, 0, nil)
}
