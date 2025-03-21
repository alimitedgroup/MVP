package application

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
)

type applyCatalogUpdatePortMock struct {
	M map[string]model.GoodInfo
}

func newApplyCatalogUpdatePortMock() *applyCatalogUpdatePortMock {
	return &applyCatalogUpdatePortMock{M: make(map[string]model.GoodInfo)}
}

func (m *applyCatalogUpdatePortMock) ApplyCatalogUpdate(good model.GoodInfo) error {
	m.M[good.ID] = good
	return nil
}

func TestApplyCatalogUpdateService(t *testing.T) {
	ctx := t.Context()

	mock := newApplyCatalogUpdatePortMock()

	app := fx.New(
		fx.Supply(fx.Annotate(mock, fx.As(new(port.IApplyCatalogUpdatePort)))),
		fx.Provide(fx.Annotate(NewApplyCatalogUpdateService, fx.As(new(port.IApplyCatalogUpdateUseCase)))),
		fx.Invoke(func(useCase port.IApplyCatalogUpdateUseCase) {
			cmd := port.CatalogUpdateCmd{
				GoodId:      "1",
				Name:        "hat",
				Description: "very nice hat",
			}

			err := useCase.ApplyCatalogUpdate(ctx, cmd)
			if err != nil {
				t.Errorf("error updating catalog: %v", err)
			}

			assert.Equal(t, mock.M["1"].Name, "hat")
		}),
	)

	err := app.Start(ctx)
	if err != nil {
		t.Errorf("error starting app: %v", err)
	}

	defer func() {
		err := app.Stop(ctx)
		if err != nil {
			t.Errorf("error stopping app: %v", err)
		}
	}()

}
