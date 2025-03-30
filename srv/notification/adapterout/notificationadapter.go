package adapterout

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/influxdata/influxdb-client-go/v2"
)

type NotificationAdapter struct {
	influxClient influxdb2.Client
	influxOrg    string
	influxBucket string
	brk          *broker.NatsMessageBroker
	ruleRepo     portout.RuleRepository
}

func NewNotificationAdapter(influxClient influxdb2.Client, brk *broker.NatsMessageBroker, ruleRepo portout.RuleRepository) *NotificationAdapter {
	return &NotificationAdapter{
		influxClient: influxClient,
		influxOrg:    "my-org",
		influxBucket: "stockdb",
		brk:          brk,
		ruleRepo:     ruleRepo,
	}
}

// =========== StockRepository port-out ===========

func (na *NotificationAdapter) SaveStockUpdate(cmd *types.AddStockUpdateCmd) error {
	writeAPI := na.influxClient.WriteAPIBlocking(na.influxOrg, na.influxBucket)
	if len(cmd.Goods) == 0 {
		return errors.New("no goods provided")
	}
	good := cmd.Goods[0]
	p := influxdb2.NewPoint(
		"stock_measurement",
		map[string]string{"warehouse_id": cmd.WarehouseID, "good_id": good.ID},
		map[string]interface{}{"quantity": good.Quantity, "delta": good.Delta, "type": cmd.Type},
		time.Now(),
	)
	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		log.Printf("Error saving to InfluxDB: %v", err)
		return err
	}
	return nil
}

// =========== StockEventPublisher port-out ===========

func (na *NotificationAdapter) PublishStockAlert(alert types.StockAlertEvent) error {
	data, err := json.Marshal(alert)
	if err != nil {
		log.Printf("Error marshalling alert: %v", err)
		return err
	}
	if err := na.brk.Nats.Publish("stock.alert", data); err != nil {
		log.Printf("Error publishing alert to NATS: %v", err)
		return err
	}
	return nil
}

// =========== RuleQueryRepository port-out ===========

func (na *NotificationAdapter) GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse {
	queryAPI := na.influxClient.QueryAPI(na.influxOrg)
	fluxQuery := `from(bucket:"stockdb")
		|> range(start:-7d)
		|> filter(fn:(r)=> r["_measurement"]=="stock_measurement")
		|> filter(fn:(r)=> r["good_id"]=="` + goodID + `")
		|> filter(fn:(r)=> r["_field"]=="quantity")
		|> last()`
	result, err := queryAPI.Query(context.Background(), fluxQuery)
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
