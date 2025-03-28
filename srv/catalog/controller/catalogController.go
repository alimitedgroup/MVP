package controller

import (
	"context"
	"encoding/json"
	"sync"

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
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	TotalRequestCounter        metric.Int64Counter
	GoodsRequestCounter        metric.Int64Counter
	GoodsGlobalQuantityCounter metric.Int64Counter
	WarehouseRequestCounter    metric.Int64Counter
	SetGoodDataCounter         metric.Int64Counter
	SetGoodQuantityCounter     metric.Int64Counter
	Logger                     *zap.Logger
	metricMap                  sync.Map
)

type catalogController struct {
	getGoodsInfoUseCase             serviceportin.IGetGoodsInfoUseCase
	getGoodsQuantityUseCase         serviceportin.IGetGoodsQuantityUseCase
	getWarehouseInfoUseCase         serviceportin.IGetWarehousesUseCase
	setMultipleGoodsQuantityUseCase serviceportin.ISetMultipleGoodsQuantityUseCase
	updateGoodDataUseCase           serviceportin.IUpdateGoodDataUseCase
}

type CatalogControllerParams struct {
	fx.In
	GetGoodsInfoUseCase             serviceportin.IGetGoodsInfoUseCase
	GetGoodsQuantityUseCase         serviceportin.IGetGoodsQuantityUseCase
	GetWarehouseInfoUseCase         serviceportin.IGetWarehousesUseCase
	SetMultipleGoodsQuantityUseCase serviceportin.ISetMultipleGoodsQuantityUseCase
	UpdateGoodDataUseCase           serviceportin.IUpdateGoodDataUseCase
	Logger                          *zap.Logger
	Meter                           metric.Meter
}

func NewCatalogController(p CatalogControllerParams) *catalogController {
	observability.CounterSetup(&p.Meter, p.Logger, &TotalRequestCounter, &metricMap, "num_catalog_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &GoodsRequestCounter, &metricMap, "num_goods_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &GoodsGlobalQuantityCounter, &metricMap, "num_goods_quantity_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &WarehouseRequestCounter, &metricMap, "num_warehouse_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &SetGoodDataCounter, &metricMap, "num_good_data_requests")
	observability.CounterSetup(&p.Meter, p.Logger, &SetGoodQuantityCounter, &metricMap, "num_good_quantity_requests")
	Logger = p.Logger
	return &catalogController{getGoodsInfoUseCase: p.GetGoodsInfoUseCase, getGoodsQuantityUseCase: p.GetGoodsQuantityUseCase, getWarehouseInfoUseCase: p.GetWarehouseInfoUseCase, setMultipleGoodsQuantityUseCase: p.SetMultipleGoodsQuantityUseCase, updateGoodDataUseCase: p.UpdateGoodDataUseCase}
}

func (cc *catalogController) getGoodsRequest(ctx context.Context, msg *nats.Msg) error { //GetGoodsInfo
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

func (cc *catalogController) getWarehouseRequest(ctx context.Context, msg *nats.Msg) error { //GetWarehouses

	Logger.Info("Received getWarehouse Request")
	verdict := "success"
	defer func() {
		Logger.Info("Completed getWarehouse request")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		WarehouseRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	req := &request.GetWarehousesInfoDTO{}

	err := json.Unmarshal(msg.Data, req)

	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		err = broker.RespondToMsg(msg, dto.GetWarehouseResponseDTO{WarehouseMap: make(map[string]dto.Warehouse), Err: err.Error()})
		if err != nil {
			Logger.Debug("Cannot send response", zap.Error(err))
		}
		return err
	}

	responseFromService := cc.getWarehouseInfoUseCase.GetWarehouses(servicecmd.NewGetWarehousesCmd())

	err = broker.RespondToMsg(msg, dto.GetWarehouseResponseDTO{WarehouseMap: responseFromService.GetWarehouseMap(), Err: ""})

	if err != nil {
		Logger.Debug("Cannot send response", zap.Error(err))
		return err
	}

	return nil
}

func (cc *catalogController) getGoodsGlobalQuantityRequest(ctx context.Context, msg *nats.Msg) error { //GetGoodsQuantity

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

func (cc *catalogController) checkSetGoodDataRequest(request *stream.GoodUpdateData) error {
	if request.GoodID == "" || request.GoodNewName == "" || request.GoodNewDescription == "" {
		return catalogCommon.ErrRequestNotValid
	}
	return nil
}

func (cc *catalogController) setGoodDataRequest(ctx context.Context, msg jetstream.Msg) error { //AddOrChangeGoodData

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

func (cc *catalogController) checkSetGoodQuantityRequest(request *stream.StockUpdate) error {
	if request.WarehouseID == "" || len(request.Goods) == 0 || request.Goods == nil {
		return catalogCommon.ErrRequestNotValid
	}
	return nil
}

func (cc *catalogController) setGoodQuantityRequest(ctx context.Context, msg jetstream.Msg) error { //SetMultipleGoodsQuantity

	Logger.Info("Received setGoodQuantity Request")
	verdict := "success"
	defer func() {
		Logger.Info("Completed setGoodQuantity request")
		TotalRequestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
		SetGoodQuantityCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verdict", verdict)))
	}()

	request := &stream.StockUpdate{}

	err := json.Unmarshal(msg.Data(), request)

	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	err = cc.checkSetGoodQuantityRequest(request)

	if err != nil {
		verdict = "bad request"
		Logger.Debug("Bad request", zap.Error(err))
		return err
	}

	responseFromService := cc.setMultipleGoodsQuantityUseCase.SetMultipleGoodsQuantity(servicecmd.NewSetMultipleGoodsQuantityCmd(request.WarehouseID, request.Goods))

	if responseFromService.GetOperationResult() == catalogCommon.ErrGenericFailure {
		Logger.Debug("Cannot complete operation", zap.Error(catalogCommon.ErrGenericFailure))
		return catalogCommon.ErrGenericFailure
	}

	return nil
}
