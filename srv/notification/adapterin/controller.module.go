package adapterin

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		AsJsController(NewAddQueryController),
		AsJsController(NewStockUpdateReceiver),
	),
	fx.Invoke(fx.Annotate(RegisterRoutes, fx.ParamTags(`group:"routes"`, `group:"js_routes"`))),
)

// RegisterRoutes è una funzione che, quando viene chiamata, registra tutti i controller presso il broker.NatsMessageBroker.
// Per funzionare, fa uso di una funzionalità di fx chiamata "Value Groups".
// Per ulteriori informazioni, fare riferimento a https://uber-go.github.io/fx/value-groups/index.html
func RegisterRoutes(core []Controller, js []JsController, brk *broker.NatsMessageBroker, rsc *broker.RestoreStreamControl) error {
	for _, controller := range core {
		err := brk.RegisterRequest(context.TODO(), controller.Subject(), broker.NoQueue, controller.Handle)
		if err != nil {
			return err
		}
	}

	for _, controller := range js {
		err := brk.RegisterJsHandler(context.TODO(), rsc, controller.Stream(), controller.Handle)
		if err != nil {
			return err
		}
	}

	return nil
}

// Controller rappresenta un handler di un endpoint NATS Core
type Controller interface {
	// Handle è la funzione che verrà chiamata quando arriva un messaggio
	Handle(ctx context.Context, msg *nats.Msg) error
	// Subject ritorna il subject di questo controller
	Subject() broker.Subject
}

// AsController marca il suo parametro come un istanza di Controller,
// in modo che RegisterRoutes possa raccogliere tutte le occorrenze
func AsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(JsController)),
		fx.ResultTags(`group:"routes"`),
	)
}

// JsController rappresenta un handler di un endpoint JetStream
type JsController interface {
	// Handle è la funzione che verrà chiamata quando arriva un messaggio
	Handle(ctx context.Context, msg jetstream.Msg) error
	// Stream ritorna la configurazione dello stream su cui si vuole ascoltare
	Stream() jetstream.StreamConfig
}

// AsJsController marca il suo parametro come un istanza di JsController,
// in modo che RegisterRoutes possa raccogliere tutte le occorrenze
func AsJsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(JsController)),
		fx.ResultTags(`group:"js_routes"`),
	)
}
