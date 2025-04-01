package adapterout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"time"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/notification/config"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type AdapterParams struct {
	fx.In
	*zap.Logger

	InfluxClient influxdb2.Client
	RuleRepo     portout.RuleRepository
	Broker       *broker.NatsMessageBroker
	Config       *config.NotificationConfig
}

func NewNotificationAdapter(p AdapterParams) (*NotificationAdapter, error) {
	str, err := p.Broker.Js.CreateStream(context.Background(), stream.AlertConfig)
	if err != nil {
		p.Logger.Error("Error creating stream", zap.String("stream", stream.AlertConfig.Name), zap.Error(err))
		return nil, err
	}

	return &NotificationAdapter{
		influx:   p.InfluxClient,
		brk:      p.Broker,
		ruleRepo: p.RuleRepo,
		writeApi: p.InfluxClient.WriteAPI(p.Config.InfluxOrg, p.Config.InfluxBucket),
		queryApi: p.InfluxClient.QueryAPI(p.Config.InfluxOrg),
		cfg:      p.Config,
		Logger:   p.Logger,
		str:      str,
	}, nil
}

type NotificationAdapter struct {
	ruleRepo portout.RuleRepository
	writeApi api.WriteAPI
	queryApi api.QueryAPI

	influx influxdb2.Client
	brk    *broker.NatsMessageBroker
	cfg    *config.NotificationConfig
	str    jetstream.Stream
	*zap.Logger
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
			map[string]interface{}{"quantity": resp.GoodMap[good.ID]},
			time.Now(),
		))
	}

	return nil
}

// =========== StockEventPublisher port-out ===========

func (na *NotificationAdapter) PublishStockAlert(alert types.StockAlertEvent) error {
	subj := fmt.Sprintf("stock.alert.%s.%s", na.cfg.ServiceId, alert.RuleId)

	for {
		var opts []jetstream.PublishOpt

		var apiErr *jetstream.APIError
		msg, err := na.str.GetLastMsgForSubject(context.Background(), subj)
		if errors.As(err, &apiErr) && apiErr.ErrorCode == jetstream.JSErrCodeMessageNotFound {
			na.Debug("No previous alert", zap.String("ruleId", alert.RuleId), zap.Error(err))
		} else if err != nil {
			na.Error("Error fetching stock alert", zap.String("subject", subj), zap.Error(err))
			return err
		} else {
			var deserialized types.StockAlertEvent
			err := json.Unmarshal(msg.Data, &deserialized)
			if err != nil {
				na.Error("Error deserializing stock alert", zap.Error(err))
				return err
			}

			if deserialized.Status != types.StockAcknowledged {
				opts = append(opts, jetstream.WithExpectLastSequence(msg.Sequence))
				_ = opts
			} else {
				na.Debug(
					"Skipping alert",
					zap.String("rule_id", alert.RuleId),
					zap.String("status", string(deserialized.Status)),
				)
				return nil
			}
		}

		data, err := json.Marshal(alert)
		if err != nil {
			na.Error("Error marshalling alert", zap.Error(err))
			return err
		}

		err = na.brk.Nats.PublishMsg(&nats.Msg{Subject: subj, Data: data})
		if errors.As(err, &apiErr) && apiErr.ErrorCode == jetstream.JSErrCodeStreamWrongLastSequence {
			na.Debug("Retrying submission, wrong sequence found")
			continue
		} else if err != nil {
			na.Error(
				"Error publishing alert to NATS",
				zap.String("ruleId", alert.RuleId),
				zap.String("status", string(alert.Status)),
				zap.String("subject", subj),
				zap.Error(err),
			)
			return err
		} else {
			return nil
		}
	}
}

func (na *NotificationAdapter) RevokeStockAlert(alert types.StockAlertEvent) error {
	for {
		var opts []jetstream.PublishOpt

		var apiErr *jetstream.APIError
		msg, err := na.str.GetLastMsgForSubject(context.Background(), fmt.Sprintf("stock.alert.%s.%s", na.cfg.ServiceId, alert.RuleId))
		if errors.As(err, &apiErr) && apiErr.ErrorCode == jetstream.JSErrCodeMessageNotFound {
			na.Debug("Skipping alert, since no previous alert was found", zap.String("rule_id", alert.RuleId))
			return nil
		} else if err != nil {
			na.Error("Error fetching stock alert", zap.Error(err))
			return err
		} else {
			var deserialized types.StockAlertEvent
			err := json.Unmarshal(msg.Data, &deserialized)
			if err != nil {
				na.Error("Error deserializing stock alert", zap.Error(err))
				return err
			}

			if deserialized.Status != types.StockRevoked {
				opts = append(opts, jetstream.WithExpectLastSequence(msg.Sequence))
			} else {
				na.Debug(
					"Skipping alert",
					zap.String("rule_id", alert.RuleId),
					zap.String("status", string(deserialized.Status)),
				)
				return nil
			}
		}

		data, err := json.Marshal(alert)
		if err != nil {
			na.Error("Error marshalling alert", zap.Error(err))
			return err
		}

		_, err = na.brk.Js.Publish(
			context.Background(),
			stream.AlertConfig.Name,
			data,
			opts...,
		)
		if errors.As(err, &apiErr) && apiErr.ErrorCode == jetstream.JSErrCodeStreamWrongLastSequence {
			continue
		} else if err != nil {
			na.Error("Error publishing alert to NATS", zap.Error(err))
			return err
		} else {
			return nil
		}
	}
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
