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

func TestReservationController(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	mock := NewMockICreateReservationUseCase(ctrl)
	mock.EXPECT().CreateReservation(gomock.Any(), gomock.Any()).Return(port.CreateReservationResponse{ReservationID: "1"}, nil)

	ns, _ := broker.NewInProcessNATSServer(t)
	cfg := config.WarehouseConfig{ID: "1"}

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t, &cfg),
		fx.Supply(fx.Annotate(mock, fx.As(new(port.ICreateReservationUseCase)))),
		fx.Provide(NewReservationController),
		fx.Provide(NewReservationRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *ReservationRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					dto := request.ReserveStockRequestDTO{
						Goods: []request.ReserveStockItem{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}
					payload, err := json.Marshal(dto)
					require.NoError(t, err)

					resp, err := ns.Request(fmt.Sprintf("warehouse.%s.reservation.create", cfg.ID), payload, 1*time.Second)
					require.NoError(t, err)

					var respDto response.ReserveStockResponseDTO
					err = json.Unmarshal(resp.Data, &respDto)
					require.NoError(t, err)

					require.Empty(t, respDto.Error)
					require.Equal(t, respDto.Message, response.ReserveStockInfo{ReservationID: "1"})

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
