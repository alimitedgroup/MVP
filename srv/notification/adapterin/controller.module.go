package adapterin

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		AsJsController(NewAddQueryController),
		AsJsController(NewStockUpdateReceiver),
	),
	fx.Invoke(fx.Annotate(RegisterRoutes, fx.ParamTags("", "", `group:"js_routes"`))),
)

// JsController rappresenta un handler di un endpoint JetStream
type JsController interface {
	// Handle è la funzione che verrà chiamata quando arriva un messaggio
	Handle(ctx context.Context, msg jetstream.Msg) error
	// Stream ritorna la configurazione dello stream su cui si vuole ascoltare
	Stream() jetstream.StreamConfig
}

// RegisterRoutes è una funzione che, quando viene chiamata, registra tutti i controller presso il broker.NatsMessageBroker.
// Per funzionare, fa uso di una funzionalità di fx chiamata "Value Groups".
// Per ulteriori informazioni, fare riferimento a https://uber-go.github.io/fx/value-groups/index.html
func RegisterRoutes(brk *broker.NatsMessageBroker, rsc *broker.RestoreStreamControl, controllers []JsController) error {
	for _, controller := range controllers {
		err := brk.RegisterJsHandler(context.TODO(), rsc, controller.Stream(), controller.Handle)
		if err != nil {
			return err
		}
	}

	return nil
}

// AsJsController marca il suo parametro come un istanza di JsController,
// in modo che RegisterRoutes possa raccoglierne tutte le istanze
func AsJsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(JsController)),
		fx.ResultTags(`group:"js_routes"`),
	)
}
