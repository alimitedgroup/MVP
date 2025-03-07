package application_test

import (
	"context"
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/application"
	"github.com/alimitedgroup/MVP/srv/warehouse/application/port"
	"github.com/alimitedgroup/MVP/srv/warehouse/model"
	"go.uber.org/fx"
)

type mockPortsImpl struct {
}

func newMockPortsImpl() *mockPortsImpl {
	return &mockPortsImpl{}
}

func (m *mockPortsImpl) GetStock(id string) int64 {
	return 0
}

func (m *mockPortsImpl) GetGood(id string) *model.GoodInfo {
	return nil
}

func (m *mockPortsImpl) CreateStockUpdate(ctx context.Context, cmd port.CreateStockUpdateCmd) error {
	return nil
}

func TestManageStockService(t *testing.T) {
	ctx := t.Context()
	mock := newMockPortsImpl()

	app := fx.New(
		fx.Supply(fx.Annotate(mock, fx.As(new(port.CreateStockUpdatePort)), fx.As(new(port.GetStockPort)), fx.As(new(port.GetGoodPort)))),
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
