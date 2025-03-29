package adapterout

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	serviceportout2 "github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type NotificationAdapter struct {
	influxClient influxdb2.Client
	influxOrg    string
	influxBucket string
	brk          *broker.NatsMessageBroker
	ruleRepo     serviceportout2.RuleRepository
}

func NewNotificationAdapter(influxClient influxdb2.Client, brk *broker.NatsMessageBroker, ruleRepo serviceportout2.RuleRepository) *NotificationAdapter {
	return &NotificationAdapter{
		influxClient: influxClient,
		influxOrg:    "my-org",
		influxBucket: "stockdb",
		brk:          brk,
		ruleRepo:     ruleRepo,
	}
}

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

func (na *NotificationAdapter) GetCurrentQuantityByGoodID(goodID string) *types.GetRuleResultResponse {
	queryAPI := na.influxClient.QueryAPI(na.influxOrg)
	fluxQuery := `
		from(bucket:"stockdb")
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

func (na *NotificationAdapter) AddRule(cmd types.QueryRule) (uuid.UUID, error) {
	return na.ruleRepo.AddRule(cmd)
}

func (na *NotificationAdapter) ListRules() ([]types.QueryRuleWithId, error) {
	return na.ruleRepo.ListRules()
}
