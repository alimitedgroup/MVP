package controller

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/catalog/service/portin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	SetGoodDataCounter  metric.Int64Counter
	GoodsRequestCounter metric.Int64Counter
)

type CatalogGoodInfoController struct {
	getGoodsInfoUseCase   serviceportin.IGetGoodsInfoUseCase
	updateGoodDataUseCase serviceportin.IUpdateGoodDataUseCase
}

func NewCatalogGoodInfoController(p CatalogControllerParams) *CatalogGoodInfoController {
	observability.CounterSetup(&p.Meter, p.Logger, &GoodsRequestCounter, &metricMap, "num_goods_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &TotalRequestCounter, &metricMap, "num_catalog_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &SetGoodDataCounter, &metricMap, "num_good_data_requests")
	Logger = p.Logger
	return &CatalogGoodInfoController{getGoodsInfoUseCase: p.GetGoodsInfoUseCase, updateGoodDataUseCase: p.UpdateGoodDataUseCase}
}

func (cc *CatalogGoodInfoController) SetGoodDataRequest(ctx context.Context, msg jetstream.Msg) error { //AddOrChangeGoodData

	Logger.Info("Received setGoodData Request")
	verdict := "success"
	defer func() {
		Logger.Info("Completed setGoodData request")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		SetGoodDataCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	request := &stream.GoodUpdateData{}

	err := json.Unmarshal(msg.Data(), request)

	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	err = cc.checkSetGoodDataRequest(request)

	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	responseFromService := cc.updateGoodDataUseCase.AddOrChangeGoodData(servicecmd.NewAddChangeGoodCmd(request.GoodID, request.GoodNewName, request.GoodNewDescription))

	if responseFromService.GetOperationResult() == catalogCommon.ErrGenericFailure {
		Logger.Debug("Cannot complete operation", zap.Error(catalogCommon.ErrGenericFailure))
		return catalogCommon.ErrGenericFailure
	}

	return nil
}

func (cc *CatalogGoodInfoController) checkSetGoodDataRequest(request *stream.GoodUpdateData) error {
	if request.GoodID == "" || request.GoodNewName == "" || request.GoodNewDescription == "" {
		return catalogCommon.ErrRequestNotValid
	}
	return nil
}

func (cc *CatalogGoodInfoController) GetGoodsRequest(ctx context.Context, msg *nats.Msg) error { //GetGoodsInfo
	Logger.Info("Received getGoods Request")
	verdict := "success"
	defer func() {
		Logger.Info("Completed getGoods request")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		GoodsRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	req := &request.GetGoodsInfoDTO{}

	err := json.Unmarshal(msg.Data, req)

	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		err = broker.RespondToMsg(msg, dto.GetGoodsDataResponseDTO{GoodMap: make(map[string]dto.Good), Err: err.Error()})
		if err != nil {
			Logger.Debug("Cannot send response", zap.Error(err))
		}
		return err
	}

	responseFromService := cc.getGoodsInfoUseCase.GetGoodsInfo(servicecmd.NewGetGoodsInfoCmd())

	err = broker.RespondToMsg(msg, dto.GetGoodsDataResponseDTO{GoodMap: responseFromService.GetMap(), Err: ""})
	if err != nil {
		Logger.Debug("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}
