package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
	service_portIn "github.com/alimitedgroup/MVP/srv/catalog/service/portIn"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

// INIZIO MOCK PORTE CONTROLLER

type FakeControllerUC struct {
	A bool
}

var (
	test *FakeControllerUC
	once sync.Once
)

func NewFakeControllerUC() *FakeControllerUC {
	once.Do(func() {
		test = &FakeControllerUC{}
	})
	return test
}

func (f *FakeControllerUC) AddOrChangeGoodData(agc *service_Cmd.AddChangeGoodCmd) *service_Response.AddOrChangeResponse {
	f.A = true
	if agc.GetId() == "wrong-test-ID" {
		return service_Response.NewAddOrChangeResponse("Not a valid goodID")
	}
	return service_Response.NewAddOrChangeResponse("Success")
}

func (f *FakeControllerUC) SetMultipleGoodsQuantity(cmd *service_Cmd.SetMultipleGoodsQuantityCmd) *service_Response.SetMultipleGoodsQuantityResponse {
	fmt.Println("LOL ", f.A)
	f.A = true
	fmt.Println("LOL ", f.A)
	errorSlice := []int{}
	for i := range cmd.GetGoods() {
		if cmd.GetGoods()[i].GoodID == "wrong-test-ID" {
			errorSlice = append(errorSlice, i)
		}
	}

	if len(errorSlice) == 0 {
		return service_Response.NewSetMultipleGoodsQuantityResponse("Success", []string{})
	}

	returnErrorSlice := []string{}
	for range errorSlice {
		returnErrorSlice = append(returnErrorSlice, "wrong-test-ID")
	}

	return service_Response.NewSetMultipleGoodsQuantityResponse("Errors", returnErrorSlice)

}

func (f *FakeControllerUC) GetGoodsQuantity(ggqc *service_Cmd.GetGoodsQuantityCmd) *service_Response.GetGoodsQuantityResponse {
	goodMap := map[string]int64{}
	goodMap["test-ID"] = int64(7)
	return service_Response.NewGetGoodsQuantityResponse(goodMap)
}

func (f *FakeControllerUC) GetGoodsInfo(ggqc *service_Cmd.GetGoodsInfoCmd) *service_Response.GetGoodsInfoResponse {
	goods := make(map[string]catalogCommon.Good)
	goods["test-ID"] = *catalogCommon.NewGood("test-ID", "test-name", "test-description")
	return service_Response.NewGetGoodsInfoResponse(goods)
}

func (f *FakeControllerUC) GetWarehouses(gwc *service_Cmd.GetWarehousesCmd) *service_Response.GetWarehousesResponse {
	warehouses := make(map[string]catalogCommon.Warehouse)
	warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")
	return service_Response.NewGetWarehousesResponse(warehouses)
}

// FINE MOCK PORTE CONTROLLER

func TestSetMultipleGoodQuantityRequest(t *testing.T) {

	ctx := t.Context()

	ns := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(service_portIn.IGetGoodsInfoUseCase)),
				fx.As(new(service_portIn.IGetGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IGetWarehousesUseCase)),
				fx.As(new(service_portIn.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IUpdateGoodDataUseCase)),
			),
		),
		fx.Provide(NewFakeControllerUC),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewCatalogController),
		fx.Provide(NewCatalogRouter),
		fx.Provide(NewControllerRouter),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Invoke(func(lc fx.Lifecycle, r *catalogRouter, f *FakeControllerUC) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					f.A = false

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
					if f.A == false {
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

	ns := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(service_portIn.IGetGoodsInfoUseCase)),
				fx.As(new(service_portIn.IGetGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IGetWarehousesUseCase)),
				fx.As(new(service_portIn.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IUpdateGoodDataUseCase)),
			),
		),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewCatalogController),
		fx.Provide(NewCatalogRouter),
		fx.Provide(NewControllerRouter),
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Provide(NewFakeControllerUC),
		fx.Invoke(func(lc fx.Lifecycle, r *catalogRouter, f *FakeControllerUC) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					f.A = false

					var request = &stream.GoodUpdateData{GoodID: "test-ID", GoodNewName: "test-name", GoodNewDescription: "test-description"}

					data, err := json.Marshal(request)

					if err != nil {
						return err
					}

					js, err := jetstream.New(ns)

					if err != nil {
						return err
					}

					_, err = js.Publish(ctx, "good.update.test", data)

					if err != nil {
						return err
					}
					time.Sleep(1 * time.Second)
					if f.A == false {
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

	ns := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(service_portIn.IGetGoodsInfoUseCase)),
				fx.As(new(service_portIn.IGetGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IGetWarehousesUseCase)),
				fx.As(new(service_portIn.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IUpdateGoodDataUseCase)),
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

					var responseDTO = &response.GetGoodsDataResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					good := make(map[string]catalogCommon.Good)
					good["test-ID"] = *catalogCommon.NewGood("test-ID", "test-name", "test-description")

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

	ns := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(service_portIn.IGetGoodsInfoUseCase)),
				fx.As(new(service_portIn.IGetGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IGetWarehousesUseCase)),
				fx.As(new(service_portIn.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IUpdateGoodDataUseCase)),
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

					var responseDTO = &response.GetWarehouseResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					warehouses := make(map[string]catalogCommon.Warehouse)
					warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")

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

	ns := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns),
		fx.Provide(
			fx.Annotate(NewFakeControllerUC,
				fx.As(new(service_portIn.IGetGoodsInfoUseCase)),
				fx.As(new(service_portIn.IGetGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IGetWarehousesUseCase)),
				fx.As(new(service_portIn.ISetMultipleGoodsQuantityUseCase)),
				fx.As(new(service_portIn.IUpdateGoodDataUseCase)),
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

					var responseDTO = &response.GetGoodsQuantityResponseDTO{}

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
