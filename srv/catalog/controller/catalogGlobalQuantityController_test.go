package controller

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestGetGoodsGlobalQuantityWrongRequest(t *testing.T) {
	ctx := t.Context()
	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		fx.Supply(ns, t),
		modules,
		fx.Invoke(func(lc fx.Lifecycle, r *catalogRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}

					responseFromController, err := ns.Request("catalog.getGoodsGlobalQuantity", []byte{}, 2*time.Second)

					if err != nil {
						return err
					}

					var responseDTO = &dto.GetGoodsQuantityResponseDTO{}

					err = json.Unmarshal(responseFromController.Data, responseDTO)

					if err != nil {
						t.Error(err)
					}

					require.NotEmpty(t, responseDTO.Err)
					assert.Equal(t, make(map[string]int64), responseDTO.GoodMap)

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
		modules,
		fx.Supply(ns, t),
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
