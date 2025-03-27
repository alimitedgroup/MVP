package notificationAdapter

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	servicecmd "github.com/alimitedgroup/MVP/srv/notification/service/cmd"
	serviceportout "github.com/alimitedgroup/MVP/srv/notification/service/portout"
	serviceresponse "github.com/alimitedgroup/MVP/srv/notification/service/response"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nats-io/nats.go"
)

type NotificationAdapter struct {
	influxClient influxdb2.Client
	influxOrg    string
	influxBucket string
	natsConn     *nats.Conn
	ruleRepo     serviceportout.IRuleRepository
}

func NewNotificationAdapter(influxClient influxdb2.Client, natsConn *nats.Conn, ruleRepo serviceportout.IRuleRepository) *NotificationAdapter {
	return &NotificationAdapter{
		influxClient: influxClient,
		influxOrg:    "my-org",
		influxBucket: "stockdb",
		natsConn:     natsConn,
		ruleRepo:     ruleRepo,
	}
}

func (na *NotificationAdapter) SaveStockUpdate(cmd *servicecmd.AddStockUpdateCmd) *serviceresponse.AddStockUpdateResponse {
	writeAPI := na.influxClient.WriteAPIBlocking(na.influxOrg, na.influxBucket)
	if len(cmd.Goods) == 0 {
		return serviceresponse.NewAddStockUpdateResponse(errors.New("no goods provided"))
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
		return serviceresponse.NewAddStockUpdateResponse(err)
	}
	return serviceresponse.NewAddStockUpdateResponse(nil)
}

func (na *NotificationAdapter) PublishStockAlert(alert serviceportout.StockAlertEvent) error {
	data, err := json.Marshal(alert)
	if err != nil {
		log.Printf("Error marshalling alert: %v", err)
		return err
	}
	if err := na.natsConn.Publish("stock.alert", data); err != nil {
		log.Printf("Error publishing alert to NATS: %v", err)
		return err
	}
	return nil
}

func (na *NotificationAdapter) GetCurrentQuantityByGoodID(goodID string) *serviceresponse.GetRuleResultResponse {
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
		return serviceresponse.NewGetRuleResultResponse(goodID, 0, err)
	}
	for result.Next() {
		val, ok := result.Record().Value().(int64)
		if !ok {
			continue
		}
		return serviceresponse.NewGetRuleResultResponse(goodID, int(val), nil)
	}
	if result.Err() != nil {
		return serviceresponse.NewGetRuleResultResponse(goodID, 0, result.Err())
	}
	return serviceresponse.NewGetRuleResultResponse(goodID, 0, nil)
}

func (na *NotificationAdapter) AddRule(cmd *servicecmd.AddQueryRuleCmd) error {
	return na.ruleRepo.AddRule(cmd)
}

func (na *NotificationAdapter) GetAllRules() []servicecmd.AddQueryRuleCmd {
	return na.ruleRepo.GetAllRules()
}
