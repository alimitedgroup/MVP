package controller

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	commonobj "github.com/alimitedgroup/MVP/common"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	common "github.com/alimitedgroup/MVP/srv/authenticator/authCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/authenticator/service/cmd"
	serviceportin "github.com/alimitedgroup/MVP/srv/authenticator/service/portIn"
	serviceresponse "github.com/alimitedgroup/MVP/srv/authenticator/service/response"
	"github.com/magiconair/properties/assert"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
	"go.uber.org/zap/zaptest"
)

//INIZIO MOCK SERVICE

type fakeService struct {
}

func NewFakeService() *fakeService {
	return &fakeService{}
}

func (fs *fakeService) GetToken(cmd *servicecmd.GetTokenCmd) *serviceresponse.GetTokenResponse {
	if cmd.GetUsername() == "wrong-username" {
		return serviceresponse.NewGetTokenResponse("", common.ErrNoToken)
	} else {
		return serviceresponse.NewGetTokenResponse("test-token", nil)
	}
}

var p = fx.Options(
	fx.Provide(broker.NewNatsMessageBroker),
	fx.Provide(NewAuthRouter),
	fx.Provide(NewAuthRouterMessageBroker),
	fx.Provide(NewControllerRouter),
	fx.Provide(func() metric.Meter {
		provider := sdkmetric.NewMeterProvider()
		a := provider.Meter("test-meter")
		return a
	}),
	fx.Provide(
		fx.Annotate(NewFakeService,
			fx.As(new(serviceportin.IGetTokenUseCase)),
		)))

// FINE MOCK SERVICE

func TestGetToken(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Supply(ns),
		p,
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(NewAuthController),
		fx.Invoke(func(lc fx.Lifecycle, r *AuthRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					msg, err := json.Marshal(commonobj.AuthLoginRequest{Username: "test-username"})
					assert.Equal(t, err, nil)
					response, err2 := ns.Request("auth.login", msg, 2*time.Second)
					assert.Equal(t, err2, nil)
					var dto commonobj.AuthLoginResponse
					err3 := json.Unmarshal(response.Data, &dto)
					assert.Equal(t, err3, nil)
					assert.Equal(t, dto.Token, "test-token")
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

func TestGetTokenWithWrongUser(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Supply(ns),
		p,
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(NewAuthController),
		fx.Invoke(func(lc fx.Lifecycle, r *AuthRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					msg, err := json.Marshal(commonobj.AuthLoginRequest{Username: "wrong-username"})
					assert.Equal(t, err, nil)
					response, err2 := ns.Request("auth.login", msg, 2*time.Second)
					assert.Equal(t, err2, nil)
					var dto commonobj.AuthLoginResponse
					err3 := json.Unmarshal(response.Data, &dto)
					assert.Equal(t, err3, nil)
					assert.Equal(t, dto.Token, "")
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

func TestGetTokenEmptyUsername(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Supply(ns),
		p,
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(NewAuthController),
		fx.Invoke(func(lc fx.Lifecycle, r *AuthRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					msg, err := json.Marshal(commonobj.AuthLoginRequest{Username: ""})
					assert.Equal(t, err, nil)
					response, err2 := ns.Request("auth.login", msg, 2*time.Second)
					assert.Equal(t, err2, nil)
					var dto commonobj.AuthLoginResponse
					err3 := json.Unmarshal(response.Data, &dto)
					assert.Equal(t, err3, nil)
					assert.Equal(t, dto.Token, "")
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

func TestGetTokenWrongRequest(t *testing.T) {
	ctx := t.Context()

	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Provide(broker.NewRestoreStreamControl),
		fx.Supply(ns),
		p,
		fx.Supply(zaptest.NewLogger(t)),
		fx.Provide(NewAuthController),
		fx.Invoke(func(lc fx.Lifecycle, r *AuthRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					assert.Equal(t, err, nil)
					response, err2 := ns.Request("auth.login", []byte(`{"nome": "test", "ciao": "test"`), 2*time.Second)
					assert.Equal(t, err2, nil)
					var dto commonobj.AuthLoginResponse
					err3 := json.Unmarshal(response.Data, &dto)
					assert.Equal(t, err3, nil)
					assert.Equal(t, dto.Token, "")
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
