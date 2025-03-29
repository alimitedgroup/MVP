package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

func TestStockController(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	addStockMock := NewMockIAddStockUseCase(ctrl)
	addStockMock.EXPECT().AddStock(gomock.Any(), gomock.Any()).Return(nil)

	removeStockMock := NewMockIRemoveStockUseCase(ctrl)
	removeStockMock.EXPECT().RemoveStock(gomock.Any(), gomock.Any()).Return(nil)

	ns, _ := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t, &cfg),
		fx.Provide(NewStockController),
		fx.Provide(NewStockRouter),
		fx.Supply(fx.Annotate(addStockMock, fx.As(new(port.IAddStockUseCase)))),
		fx.Supply(fx.Annotate(removeStockMock, fx.As(new(port.IRemoveStockUseCase)))),
		fx.Invoke(func(lc fx.Lifecycle, r *StockRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					addDto := request.AddStockRequestDTO{
						GoodID:   "1",
						Quantity: 10,
					}
					addPayload, err := json.Marshal(addDto)
					require.NoError(t, err)

					addResp, err := ns.Request(fmt.Sprintf("warehouse.%s.stock.add", cfg.ID), addPayload, 1*time.Second)
					require.NoError(t, err)

					var addRespDto response.ResponseDTO[string]
					err = json.Unmarshal(addResp.Data, &addRespDto)
					require.NoError(t, err)

					require.Equal(t, addRespDto.Message, "ok")

					remDto := request.AddStockRequestDTO{
						GoodID:   "1",
						Quantity: 10,
					}
					remPayload, err := json.Marshal(remDto)
					require.NoError(t, err)

					remResp, err := ns.Request(fmt.Sprintf("warehouse.%s.stock.remove", cfg.ID), remPayload, 1*time.Second)
					require.NoError(t, err)

					var remRespDto response.ResponseDTO[string]
					err = json.Unmarshal(remResp.Data, &remRespDto)
					require.NoError(t, err)

					require.Equal(t, remRespDto.Message, "ok")

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

func TestStockControllerAddStockErr(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	addStockMock := NewMockIAddStockUseCase(ctrl)
	addStockMock.EXPECT().AddStock(gomock.Any(), gomock.Any()).Return(fmt.Errorf("mock error"))

	removeStockMock := NewMockIRemoveStockUseCase(ctrl)

	ns, _ := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t, &cfg),
		fx.Provide(NewStockController),
		fx.Provide(NewStockRouter),
		fx.Supply(fx.Annotate(addStockMock, fx.As(new(port.IAddStockUseCase)))),
		fx.Supply(fx.Annotate(removeStockMock, fx.As(new(port.IRemoveStockUseCase)))),
		fx.Invoke(func(lc fx.Lifecycle, r *StockRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					addDto := request.AddStockRequestDTO{
						GoodID:   "1",
						Quantity: 10,
					}
					addPayload, err := json.Marshal(addDto)
					require.NoError(t, err)

					addResp, err := ns.Request(fmt.Sprintf("warehouse.%s.stock.add", cfg.ID), addPayload, 1*time.Second)
					require.NoError(t, err)

					var addRespDto response.ResponseDTO[string]
					err = json.Unmarshal(addResp.Data, &addRespDto)
					require.NoError(t, err)

					require.Empty(t, addRespDto.Message)
					require.NotEmpty(t, addRespDto.Error)

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

func TestStockControllerRemStockErr(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	addStockMock := NewMockIAddStockUseCase(ctrl)

	removeStockMock := NewMockIRemoveStockUseCase(ctrl)
	removeStockMock.EXPECT().RemoveStock(gomock.Any(), gomock.Any()).Return(fmt.Errorf("mock error"))

	ns, _ := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t, &cfg),
		fx.Provide(NewStockController),
		fx.Provide(NewStockRouter),
		fx.Supply(fx.Annotate(addStockMock, fx.As(new(port.IAddStockUseCase)))),
		fx.Supply(fx.Annotate(removeStockMock, fx.As(new(port.IRemoveStockUseCase)))),
		fx.Invoke(func(lc fx.Lifecycle, r *StockRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					remDto := request.AddStockRequestDTO{
						GoodID:   "1",
						Quantity: 10,
					}
					remPayload, err := json.Marshal(remDto)
					require.NoError(t, err)

					remResp, err := ns.Request(fmt.Sprintf("warehouse.%s.stock.remove", cfg.ID), remPayload, 1*time.Second)
					require.NoError(t, err)

					var remRespDto response.ResponseDTO[string]
					err = json.Unmarshal(remResp.Data, &remRespDto)
					require.NoError(t, err)

					require.Empty(t, remRespDto.Message)
					require.NotEmpty(t, remRespDto.Error)

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
