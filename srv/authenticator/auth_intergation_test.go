package main

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	commonobj "github.com/alimitedgroup/MVP/common"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/authenticator/config"
	"github.com/alimitedgroup/MVP/srv/authenticator/controller"
	"github.com/magiconair/properties/assert"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var (
	check bool
	key   string
	mutex sync.Mutex
)

func setKey(value string) {
	mutex.Lock()
	defer mutex.Unlock()
	key = value
}

func getKey() string {
	mutex.Lock()
	value := key
	mutex.Unlock()
	return value
}

func setCheck(value bool) {
	mutex.Lock()
	defer mutex.Unlock()
	check = value
}

func getCheck() bool {
	mutex.Lock()
	value := check
	mutex.Unlock()
	return value
}

func checkKeyStream(ctx context.Context, msg jetstream.Msg) error {
	if msg.Data() != nil {
		setKey(string(msg.Data()))
		setCheck(true)
	}
	return nil
}

func executeCheck(ctx context.Context, msg jetstream.Msg) error {
	if string(msg.Data()) == getKey() {
		setKey(string(msg.Data()))
		setCheck(true)
	}
	return nil
}

func chekKey(broker *broker.NatsMessageBroker, rsc *broker.RestoreStreamControl) error {
	err := broker.RegisterJsHandler(context.Background(), rsc, stream.KeyStream, executeCheck)
	if err != nil {
		return err
	}
	return nil
}

func TestGetTokenEmptyUsername(t *testing.T) {
	ctx := t.Context()
	setCheck(false)
	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Supply(ns),
		config.Modules,
		fx.Invoke(func(lc fx.Lifecycle, r *controller.AuthRouter, broker *broker.NatsMessageBroker, rsc *broker.RestoreStreamControl) {
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
					if dto.Token != "" {
						t.Error("Expected empty Token")
					}
					ctx = context.Background()
					err = broker.RegisterJsHandler(ctx, rsc, stream.KeyStream, checkKeyStream)
					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)
					assert.Equal(t, getCheck(), false) //Not a valid username and it was the first token requested, so no key has been generated
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

func TestGetTokenWrongUsername(t *testing.T) {
	ctx := t.Context()
	setCheck(false)
	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Supply(ns),
		config.Modules,
		fx.Invoke(func(lc fx.Lifecycle, r *controller.AuthRouter, broker *broker.NatsMessageBroker, rsc *broker.RestoreStreamControl) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					if err != nil {
						t.Error(err)
					}
					msg, err := json.Marshal(commonobj.AuthLoginRequest{Username: "test"})
					assert.Equal(t, err, nil)
					response, err2 := ns.Request("auth.login", msg, 2*time.Second)
					assert.Equal(t, err2, nil)
					var dto commonobj.AuthLoginResponse
					err3 := json.Unmarshal(response.Data, &dto)
					assert.Equal(t, err3, nil)
					if dto.Token != "" {
						t.Error("Expected empty Token")
					}
					ctx = context.Background()
					err = broker.RegisterJsHandler(ctx, rsc, stream.KeyStream, checkKeyStream)
					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)
					assert.Equal(t, getCheck(), false) //Not a valid username and it was the first token requested, so no key has been generated
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

func TestGetToken(t *testing.T) {
	ctx := t.Context()
	setCheck(false)
	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Supply(ns),
		config.Modules,
		fx.Invoke(func(lc fx.Lifecycle, r *controller.AuthRouter, broker *broker.NatsMessageBroker, rsc *broker.RestoreStreamControl) {
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
					ctx = context.Background()
					err = broker.RegisterJsHandler(ctx, rsc, stream.KeyStream, checkKeyStream)
					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)
					assert.Equal(t, getCheck(), true)
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

func TestGetTwoToken(t *testing.T) {
	ctx := t.Context()
	setCheck(false)
	ns, _ := broker.NewInProcessNATSServer(t)
	app := fx.New(
		fx.Supply(ns),
		config.Modules,
		fx.Invoke(func(lc fx.Lifecycle, r *controller.AuthRouter, broker *broker.NatsMessageBroker, rsc *broker.RestoreStreamControl) {
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
					ctx = context.Background()
					err = broker.RegisterJsHandler(ctx, rsc, stream.KeyStream, checkKeyStream)
					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)
					assert.Equal(t, getCheck(), true)
					//secondo token
					setCheck(false)
					msg, err = json.Marshal(commonobj.AuthLoginRequest{Username: "client"})
					assert.Equal(t, err, nil)
					response, err2 = ns.Request("auth.login", msg, 2*time.Second)
					assert.Equal(t, err2, nil)
					var dto2 commonobj.AuthLoginResponse
					err3 = json.Unmarshal(response.Data, &dto2)
					assert.Equal(t, err3, nil)
					if dto.Token == "" {
						t.Error("Second token is empty")
					}
					time.Sleep(1 * time.Second)
					err = chekKey(broker, rsc)
					if err != nil {
						return err
					}
					time.Sleep(2 * time.Second)
					assert.Equal(t, getCheck(), true)
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
