package listener

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type StockUpdateListener struct {
}

func NewStockUpdateListener() *StockUpdateListener {
	return &StockUpdateListener{}
}

func (l *StockUpdateListener) ListenStockUpdate(ctx context.Context, msg jetstream.Msg) error {
	return nil
}
