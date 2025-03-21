package business

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

type applyCatalogUpdatePortMock struct {
	M map[string]model.GoodInfo
}

func newApplyCatalogUpdatePortMock() *applyCatalogUpdatePortMock {
	return &applyCatalogUpdatePortMock{M: make(map[string]model.GoodInfo)}
}

func (m *applyCatalogUpdatePortMock) ApplyCatalogUpdate(good model.GoodInfo) {
	m.M[good.ID] = good
}

func TestApplyCatalogUpdateService(t *testing.T) {
	ctx := t.Context()

	mock := newApplyCatalogUpdatePortMock()

	app := fx.New(
		fx.Supply(fx.Annotate(mock, fx.As(new(port.IApplyCatalogUpdatePort)))),
		fx.Provide(fx.Annotate(NewApplyCatalogUpdateService, fx.As(new(port.IApplyCatalogUpdateUseCase)))),
		fx.Invoke(func(useCase port.IApplyCatalogUpdateUseCase) {
			cmd := port.CatalogUpdateCmd{
				GoodID:      "1",
				Name:        "hat",
				Description: "very nice hat",
			}

			useCase.ApplyCatalogUpdate(cmd)

			assert.Equal(t, mock.M["1"].Name, "hat")
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)
	defer func() {
		err := app.Stop(ctx)
		require.NoError(t, err)
	}()

}
