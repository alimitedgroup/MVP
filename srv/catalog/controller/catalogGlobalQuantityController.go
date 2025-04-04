package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/catalog/service/portin"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	GoodsGlobalQuantityCounter metric.Int64Counter
)

type CatalogGlobalQuantityController struct {
	getGoodsQuantityUseCase serviceportin.IGetGoodsQuantityUseCase
}

func NewCatalogGlobalQuantityController(p CatalogControllerParams) *CatalogGlobalQuantityController {
	observability.CounterSetup(&p.Meter, p.Logger, &TotalRequestCounter, &metricMap, "num_catalog_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &GoodsGlobalQuantityCounter, &metricMap, "num_goods_quantity_requests")
	Logger = p.Logger
	return &CatalogGlobalQuantityController{getGoodsQuantityUseCase: p.GetGoodsQuantityUseCase}
}

func (cc *CatalogGlobalQuantityController) GetGoodsGlobalQuantityRequest(ctx context.Context, msg *nats.Msg) error { //GetGoodsQuantity

	Logger.Info("Received GetGoodsGlobalQuantity Request")
	verdict := "success"
	defer func() {
		Logger.Info("Completed GetGoodsGlobalQuantity request")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		GoodsGlobalQuantityCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	request := &request.GetGoodsQuantityDTO{}

	err := json.Unmarshal(msg.Data, request)

	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		Logger.Debug("Bad request", zap.Error(err))
		err = broker.RespondToMsg(msg, dto.GetGoodsQuantityResponseDTO{GoodMap: make(map[string]int64), Err: err.Error()})
		if err != nil {
			Logger.Debug("Cannot send response", zap.Error(err))
		}
		return err
	}

	responseFromService := cc.getGoodsQuantityUseCase.GetGoodsQuantity(servicecmd.NewGetGoodsQuantityCmd())

	err = broker.RespondToMsg(msg, dto.GetGoodsQuantityResponseDTO{GoodMap: responseFromService.GetMap(), Err: ""})

	if err != nil {
		Logger.Debug("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}
