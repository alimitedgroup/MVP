package application_test

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/application"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"go.uber.org/fx"
)

type GetStockPortMock struct {
}

func newGetStockPortMock() *GetStockPortMock {
	return &GetStockPortMock{}
}

func (g *GetStockPortMock) GetStock(id string) int64 {
	return 0
}

type CreateStockUpdatePortMock struct {
}

func newCreateStockUpdatePortMock() *CreateStockUpdatePortMock {
	return &CreateStockUpdatePortMock{}
}

func (c *CreateStockUpdatePortMock) CreateStockUpdate(ctx context.Context, cmd port.CreateStockUpdateCmd) error {
	return nil
}

func TestManageStockService(t *testing.T) {
	ctx := t.Context()

	app := fx.New(
		fx.Provide(newCreateStockUpdatePortMock, func(s *CreateStockUpdatePortMock) port.CreateStockUpdatePort { return s }),
		fx.Provide(newGetStockPortMock, func(s *GetStockPortMock) port.GetStockPort { return s }),

		fx.Provide(fx.Annotate(application.NewManageStockService, fx.As(new(port.AddStockUseCase)), fx.As(new(port.RemoveStockUseCase)))),
		fx.Invoke(func(addStockUseCase port.AddStockUseCase, removeStockUseCase port.RemoveStockUseCase) {}),
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
