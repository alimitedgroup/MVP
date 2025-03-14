package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	commonobj "github.com/alimitedgroup/MVP/common"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/authenticator/config"
	"github.com/alimitedgroup/MVP/srv/authenticator/controller"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestGetToken(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Supply(ns),
		config.Modules,
		fx.Invoke(func(lc fx.Lifecycle, r *controller.AuthRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					msg, err := json.Marshal(commonobj.AuthLoginRequest{Username: "client"})
					assert.Equal(t, err, nil)
					response, err2 := ns.Request("auth.login", msg, 2*time.Second)
					assert.Equal(t, err2, nil)
					var dto commonobj.AuthLoginResponse
					err3 := json.Unmarshal(response.Data, &dto)
					assert.Equal(t, err3, nil)
					if dto.Token == "" {
						t.Error("Empty Token")
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
