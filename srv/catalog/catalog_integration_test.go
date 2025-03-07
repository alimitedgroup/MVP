package main

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogAdapter"
	"github.com/alimitedgroup/MVP/srv/catalog/controller"
	goodRepository "github.com/alimitedgroup/MVP/srv/catalog/persistance"
	"github.com/alimitedgroup/MVP/srv/catalog/service"
	"github.com/magiconair/properties/assert"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var ModulesfForTesting = fx.Options(
	controller.Module,
	goodRepository.Module,
	catalogAdapter.Module,
	service.Module,
)

func TestInsertGetWarehousesQuantity(t *testing.T) {
	ns := broker.NewInProcessNATSServer(t)
	ctx := context.Background()
	app := fx.New(
		fx.Supply(ns),
		ModulesfForTesting,
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Invoke(func(lc fx.Lifecycle, r *controller.ControllerRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					//invio dati
					goods := &stream.StockUpdateGood{GoodID: "test-ID", Quantity: 7, Delta: 0}
					goods2 := &stream.StockUpdateGood{GoodID: "test-ID", Quantity: 2, Delta: 5}
					goods3 := &stream.StockUpdateGood{GoodID: "2test-ID", Quantity: 3, Delta: 0}
					goods4 := &stream.StockUpdateGood{GoodID: "test-ID", Quantity: 3, Delta: 0}
					goodsSlice := []stream.StockUpdateGood{}
					goodsSlice2 := []stream.StockUpdateGood{}
					goodsSlice = append(goodsSlice, *goods)
					goodsSlice = append(goodsSlice, *goods2)
					goodsSlice2 = append(goodsSlice2, *goods3)
					goodsSlice2 = append(goodsSlice2, *goods4)

					var request1 = &stream.StockUpdate{ID: "update-ID", WarehouseID: "test-warehouse-ID", Type: stream.StockUpdateTypeAdd, Goods: goodsSlice, OrderID: "test-order-ID", TransferID: "test-transfer-ID", Timestamp: 1}
					var request2 = &stream.StockUpdate{ID: "2update-ID", WarehouseID: "2test-warehouse-ID", Type: stream.StockUpdateTypeAdd, Goods: goodsSlice2, OrderID: "2test-order-ID", TransferID: "test-transfer-ID", Timestamp: 1}

					data, err := json.Marshal(request1)
					if err != nil {
						return err
					}
					data2, err := json.Marshal(request2)
					if err != nil {
						return err
					}

					js, err := jetstream.New(ns)
					if err != nil {
						return err
					}
					_, err = js.CreateStream(ctx, stream.StockUpdateStreamConfig)
					if err != nil {
						return err
					}
					_, err = js.Publish(ctx, "stock.update.test", data)

					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)

					_, err = js.Publish(ctx, "stock.update.test", data2)

					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)
					//recupero dati

					var request3 = &request.GetWarehousesInfoDTO{}

					data3, err := json.Marshal(request3)

					if err != nil {
						return err
					}

					responseFromController, err := ns.Request("catalog.getWarehouses", data3, 3*time.Second)

					if err != nil {
						return err
					}

					var responseDTO = &response.GetWarehouseResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, responseDTO.WarehouseMap["test-warehouse-ID"].ID, "test-warehouse-ID")
					assert.Equal(t, responseDTO.WarehouseMap["2test-warehouse-ID"].ID, "2test-warehouse-ID")
					assert.Equal(t, responseDTO.WarehouseMap["test-warehouse-ID"].Stock["test-ID"], int64(2))
					assert.Equal(t, responseDTO.WarehouseMap["2test-warehouse-ID"].Stock["test-ID"], int64(3))
					assert.Equal(t, responseDTO.WarehouseMap["test-warehouse-ID"].Stock["2test-ID"], int64(0))
					assert.Equal(t, responseDTO.WarehouseMap["2test-warehouse-ID"].Stock["2test-ID"], int64(3))
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
		log.Fatal(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func TestInsertGetGoodsQuantity(t *testing.T) {
	ns := broker.NewInProcessNATSServer(t)
	ctx := context.Background()
	app := fx.New(
		fx.Supply(ns),
		ModulesfForTesting,
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Invoke(func(lc fx.Lifecycle, r *controller.ControllerRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					//invio dati
					goods := &stream.StockUpdateGood{GoodID: "test-ID", Quantity: 7, Delta: 0}
					goods2 := &stream.StockUpdateGood{GoodID: "test-ID", Quantity: 2, Delta: 5}
					goods3 := &stream.StockUpdateGood{GoodID: "2test-ID", Quantity: 3, Delta: 0}
					goods4 := &stream.StockUpdateGood{GoodID: "test-ID", Quantity: 3, Delta: 0}
					goodsSlice := []stream.StockUpdateGood{}
					goodsSlice2 := []stream.StockUpdateGood{}
					goodsSlice = append(goodsSlice, *goods)
					goodsSlice = append(goodsSlice, *goods2)
					goodsSlice2 = append(goodsSlice2, *goods3)
					goodsSlice2 = append(goodsSlice2, *goods4)

					var request1 = &stream.StockUpdate{ID: "update-ID", WarehouseID: "test-warehouse-ID", Type: stream.StockUpdateTypeAdd, Goods: goodsSlice, OrderID: "test-order-ID", TransferID: "test-transfer-ID", Timestamp: 1}
					var request2 = &stream.StockUpdate{ID: "2update-ID", WarehouseID: "2test-warehouse-ID", Type: stream.StockUpdateTypeAdd, Goods: goodsSlice2, OrderID: "2test-order-ID", TransferID: "test-transfer-ID", Timestamp: 1}

					data, err := json.Marshal(request1)
					if err != nil {
						return err
					}
					data2, err := json.Marshal(request2)
					if err != nil {
						return err
					}

					js, err := jetstream.New(ns)
					if err != nil {
						return err
					}
					_, err = js.CreateStream(ctx, stream.StockUpdateStreamConfig)
					if err != nil {
						return err
					}
					_, err = js.Publish(ctx, "stock.update.test", data)

					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)

					_, err = js.Publish(ctx, "stock.update.test", data2)

					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)
					//recupero dati

					var request3 = &request.GetGoodsInfoDTO{}

					data3, err := json.Marshal(request3)

					if err != nil {
						return err
					}

					responseFromController, err := ns.Request("catalog.getGoodsGlobalQuantity", data3, 3*time.Second)

					if err != nil {
						return err
					}

					var responseDTO = &response.GetGoodsQuantityResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, responseDTO.GoodMap["test-ID"], int64(5))
					assert.Equal(t, responseDTO.GoodMap["2test-ID"], int64(3))

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
		log.Fatal(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func TestInsertGetGoods(t *testing.T) {
	ns := broker.NewInProcessNATSServer(t)
	ctx := context.Background()
	app := fx.New(
		fx.Supply(ns),
		ModulesfForTesting,
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Invoke(func(lc fx.Lifecycle, r *controller.ControllerRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					//Inizio invio dati

					var request1 = &stream.GoodUpdateData{GoodID: "test-ID", GoodNewName: "test-name", GoodNewDescription: "test-description"}

					data, err := json.Marshal(request1)

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
					//Inizio recupero dati
					request2 := &request.GetGoodsInfoDTO{}

					data2, err := json.Marshal(request2)

					if err != nil {
						return err
					}

					responseFromController, err := ns.Request("catalog.getGoods", data2, 2*time.Second)

					if err != nil {
						return err
					}

					var responseDTO = &response.GetGoodsDataResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					assert.Equal(t, responseDTO.GoodMap["test-ID"].ID, "test-ID")
					assert.Equal(t, responseDTO.GoodMap["test-ID"].Name, "test-name")
					assert.Equal(t, responseDTO.GoodMap["test-ID"].Description, "test-description")

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
		log.Fatal(err)
	}

	defer func() {
		err = app.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
