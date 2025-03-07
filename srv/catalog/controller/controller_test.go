package controller

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
	service_portIn "github.com/alimitedgroup/MVP/srv/catalog/service/portIn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

// INIZIO MOCK PORTE CONTROLLER

var setData bool            //if true, setDataRequest had effect
var updateDataQuantity bool //if true, setMultipleGoodQuantityRequest had effect

type fakeControllerUC struct {
}

func NewFakeControllerUC() *fakeControllerUC {
	setData = false
	updateDataQuantity = false
	return &fakeControllerUC{}
}

func (f *fakeControllerUC) AddOrChangeGoodData(agc *service_Cmd.AddChangeGoodCmd) *service_Response.AddOrChangeResponse {
	setData = true
	if agc.GetId() == "wrong-test-ID" {
		return service_Response.NewAddOrChangeResponse("Not a valid goodID")
	}
	return service_Response.NewAddOrChangeResponse("Success")
}

func (f *fakeControllerUC) SetMultipleGoodsQuantity(cmd *service_Cmd.SetMultipleGoodsQuantityCmd) *service_Response.SetMultipleGoodsQuantityResponse {
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
	updateDataQuantity = true
	for range errorSlice {
		returnErrorSlice = append(returnErrorSlice, "wrong-test-ID")
	}

	return service_Response.NewSetMultipleGoodsQuantityResponse("Errors", returnErrorSlice)

}

func (f *fakeControllerUC) GetGoodsQuantity(ggqc *service_Cmd.GetGoodsQuantityCmd) *service_Response.GetGoodsQuantityResponse {
	goodMap := map[string]int64{}
	goodMap["test-ID"] = int64(7)
	return service_Response.NewGetGoodsQuantityResponse(goodMap)
}

func (f *fakeControllerUC) GetGoodsInfo(ggqc *service_Cmd.GetGoodsInfoCmd) *service_Response.GetGoodsInfoResponse {
	goods := make(map[string]catalogCommon.Good)
	goods["test-ID"] = *catalogCommon.NewGood("test-ID", "test-name", "test-description")
	return service_Response.NewGetGoodsInfoResponse(goods)
}

func (f *fakeControllerUC) GetWarehouses(gwc *service_Cmd.GetWarehousesCmd) *service_Response.GetWarehousesResponse {
	warehouses := make(map[string]catalogCommon.Warehouse)
	warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")
	return service_Response.NewGetWarehousesResponse(warehouses)
}

// FINE MOCK PORTE CONTROLLER

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

/*func TestSetGoodDataRequest(t *testing.T) {
	ctx := t.Context()

	//ns := broker.NewInProcessNATSServer(t)
	ns, err := nats.Connect(nats.DefaultURL)

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

					var request = &stream.GoodUpdateData{GoodID: "test-ID", GoodNewName: "test-name", GoodNewDescription: "test-description"}

					data, err := json.Marshal(request)

					if err != nil {
						return err
					}

					js, err := jetstream.New(ns)

					if err != nil {
						return err
					}

					_, err = js.CreateStream(ctx, stream.AddOrChangeGoodDataStream)

					fmt.Println("KJCDSAKSAN ", err)

					_, err = js.Publish(ctx, "good.update.test", data)

					if err != nil {
						return err
					}

					assert.Equal(t, setData, true)

					setData = false

					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
				},
			})
		}),
	)

	err = app.Start(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			t.Error(err)
		}
	}()
}*/
