package controller

import (
	"context"
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/dto"
	"go.uber.org/zap/zaptest"
	"sync"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/catalog/service/portin"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

// INIZIO MOCK PORTE CONTROLLER

var (
	a     bool
	mutex sync.Mutex
)

func changeA(value bool) {
	mutex.Lock()
	a = value
	mutex.Unlock()
}

func getA() bool {
	mutex.Lock()
	value := a
	mutex.Unlock()
	return value
}

type FakeControllerUC struct {
}

func NewFakeControllerUC() *FakeControllerUC {
	return &FakeControllerUC{}
}

func (f *FakeControllerUC) AddOrChangeGoodData(agc *servicecmd.AddChangeGoodCmd) *serviceresponse.AddOrChangeResponse {
	changeA(true)
	if agc.GetId() == "wrong-test-ID" {
		return serviceresponse.NewAddOrChangeResponse(catalogCommon.ErrGoodIdNotValid)
	}
	return serviceresponse.NewAddOrChangeResponse(nil)
}

func (f *FakeControllerUC) SetMultipleGoodsQuantity(cmd *servicecmd.SetMultipleGoodsQuantityCmd) *serviceresponse.SetMultipleGoodsQuantityResponse {
	changeA(true)
	errorSlice := []int{}
	for i := range cmd.GetGoods() {
		if cmd.GetGoods()[i].GoodID == "wrong-test-ID" {
			errorSlice = append(errorSlice, i)
		}
	}

	if len(errorSlice) == 0 {
		return serviceresponse.NewSetMultipleGoodsQuantityResponse(nil, []string{})
	}

	returnErrorSlice := []string{}
	for range errorSlice {
		returnErrorSlice = append(returnErrorSlice, "wrong-test-ID")
	}

	return serviceresponse.NewSetMultipleGoodsQuantityResponse(catalogCommon.ErrGenericFailure, returnErrorSlice)

}

func (f *FakeControllerUC) GetGoodsQuantity(ggqc *servicecmd.GetGoodsQuantityCmd) *serviceresponse.GetGoodsQuantityResponse {
	goodMap := map[string]int64{}
	goodMap["test-ID"] = int64(7)
	return serviceresponse.NewGetGoodsQuantityResponse(goodMap)
}

func (f *FakeControllerUC) GetGoodsInfo(ggqc *servicecmd.GetGoodsInfoCmd) *serviceresponse.GetGoodsInfoResponse {
	goods := make(map[string]dto.Good)
	goods["test-ID"] = *dto.NewGood("test-ID", "test-name", "test-description")
	return serviceresponse.NewGetGoodsInfoResponse(goods)
}

func (f *FakeControllerUC) GetWarehouses(gwc *servicecmd.GetWarehousesCmd) *serviceresponse.GetWarehousesResponse {
	warehouses := make(map[string]dto.Warehouse)
	warehouses["test-warehouse-ID"] = *dto.NewWarehouse("test-warehose-ID")
	return serviceresponse.NewGetWarehousesResponse(warehouses)
}

// FINE MOCK PORTE CONTROLLER

func TestSetMultipleGoodQuantityRequest(t *testing.T) {

	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(serviceportin.IGetGoodsInfoUseCase)),
				fx.As(new(serviceportin.IGetGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IGetWarehousesUseCase)),
				fx.As(new(serviceportin.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IUpdateGoodDataUseCase)),
			),
		),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewCatalogController),
		fx.Provide(NewCatalogRouter),
		fx.Provide(NewControllerRouter),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Invoke(func(lc fx.Lifecycle, r *catalogRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					changeA(false)

					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					goods := &stream.StockUpdateGood{GoodID: "test-ID", Quantity: 7, Delta: 0}

					goodsSlice := []stream.StockUpdateGood{}
					goodsSlice = append(goodsSlice, *goods)

					var request = &stream.StockUpdate{ID: "update-ID", WarehouseID: "test-warehouse-ID", Type: stream.StockUpdateTypeAdd, Goods: goodsSlice, OrderID: "test-order-ID", TransferID: "test-transfer-ID", Timestamp: 1}

					data, err := json.Marshal(request)

					if err != nil {
						return err
					}

					js, err := jetstream.New(ns)

					if err != nil {
						return err
					}
					_, err = js.Publish(ctx, "stock.update.test", data)

					if err != nil {
						return err
					}
					time.Sleep(1 * time.Second)
					if getA() == false {
						t.Errorf("Expected true returned false")
					}
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)
	err := app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}
func TestSetGoodDataRequest(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(serviceportin.IGetGoodsInfoUseCase)),
				fx.As(new(serviceportin.IGetGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IGetWarehousesUseCase)),
				fx.As(new(serviceportin.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IUpdateGoodDataUseCase)),
			),
		),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewCatalogController),
		fx.Provide(NewCatalogRouter),
		fx.Provide(NewControllerRouter),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Invoke(func(lc fx.Lifecycle, r *catalogRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					changeA(false)

					var request = &stream.GoodUpdateData{GoodID: "test-ID", GoodNewName: "test-name", GoodNewDescription: "test-description"}

					data, err := json.Marshal(request)

					if err != nil {
						return err
					}

					js, err := jetstream.New(ns)

					if err != nil {
						return err
					}

					_, err = js.Publish(ctx, "good.update", data)

					if err != nil {
						return err
					}
					time.Sleep(1 * time.Second)
					if getA() == false {
						t.Errorf("Expected true returned false")
					}
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestGetGoodsRequest(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(serviceportin.IGetGoodsInfoUseCase)),
				fx.As(new(serviceportin.IGetGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IGetWarehousesUseCase)),
				fx.As(new(serviceportin.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IUpdateGoodDataUseCase)),
			),
		),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewCatalogController),
		fx.Provide(NewCatalogRouter),
		fx.Provide(NewControllerRouter),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Invoke(func(lc fx.Lifecycle, r *catalogRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					var request = &request.GetGoodsInfoDTO{}

					data, err := json.Marshal(request)

					if err != nil {
						return err
					}

					responseFromController, err := ns.Request("catalog.getGoods", data, 2*time.Second)

					if err != nil {
						return err
					}

					var responseDTO = &dto.GetGoodsDataResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					good := make(map[string]dto.Good)
					good["test-ID"] = *dto.NewGood("test-ID", "test-name", "test-description")

					assert.Equal(t, responseDTO.Err, "")
					assert.Equal(t, responseDTO.GoodMap, good)

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestGetWarehousesRequest(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(serviceportin.IGetGoodsInfoUseCase)),
				fx.As(new(serviceportin.IGetGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IGetWarehousesUseCase)),
				fx.As(new(serviceportin.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IUpdateGoodDataUseCase)),
			),
		),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewCatalogController),
		fx.Provide(NewCatalogRouter),
		fx.Provide(NewControllerRouter),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Invoke(func(lc fx.Lifecycle, r *catalogRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					var request = &request.GetWarehousesInfoDTO{}

					data, err := json.Marshal(request)

					if err != nil {
						return err
					}

					responseFromController, err := ns.Request("catalog.getWarehouses", data, 2*time.Second)

					if err != nil {
						return err
					}

					var responseDTO = &dto.GetWarehouseResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					warehouses := make(map[string]dto.Warehouse)
					warehouses["test-warehouse-ID"] = *dto.NewWarehouse("test-warehose-ID")

					assert.Equal(t, responseDTO.Err, "")
					assert.Equal(t, responseDTO.WarehouseMap, warehouses)

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestGetGoodsGlobalQuantityRequest(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(zaptest.NewLogger(t)),
		fx.Supply(ns),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(serviceportin.IGetGoodsInfoUseCase)),
				fx.As(new(serviceportin.IGetGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IGetWarehousesUseCase)),
				fx.As(new(serviceportin.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(serviceportin.IUpdateGoodDataUseCase)),
			),
		),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewCatalogController),
		fx.Provide(NewCatalogRouter),
		fx.Provide(NewControllerRouter),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Invoke(func(lc fx.Lifecycle, r *catalogRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					var request = &request.GetGoodsInfoDTO{}

					data, err := json.Marshal(request)

					if err != nil {
						return err
					}

					responseFromController, err := ns.Request("catalog.getGoodsGlobalQuantity", data, 2*time.Second)

					if err != nil {
						return err
					}

					var responseDTO = &dto.GetGoodsQuantityResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					goodMap := map[string]int64{}
					goodMap["test-ID"] = int64(7)

					assert.Equal(t, responseDTO.Err, "")
					assert.Equal(t, responseDTO.GoodMap, goodMap)

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}
