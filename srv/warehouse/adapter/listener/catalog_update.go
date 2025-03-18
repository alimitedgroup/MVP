package listener

import (
	"context"
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/nats-io/nats.go/jetstream"
)

type CatalogListener struct {
	applyCatalogUpdateUseCase port.IApplyCatalogUpdateUseCase
}

func NewCatalogListener(applyCatalogUpdateUseCase port.IApplyCatalogUpdateUseCase) *CatalogListener {
	return &CatalogListener{applyCatalogUpdateUseCase}
}

func (l *CatalogListener) ListenGoodUpdate(ctx context.Context, msg jetstream.Msg) error {
	var event stream.GoodUpdateData
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		return err
	}

	cmd := port.CatalogUpdateCmd{
		GoodID:      event.GoodID,
		Name:        event.GoodNewName,
		Description: event.GoodNewDescription,
	}

	l.applyCatalogUpdateUseCase.ApplyCatalogUpdate(cmd)

	return nil
}
