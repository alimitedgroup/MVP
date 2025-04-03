package business

import (
	"context"
	"fmt"
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/magiconair/properties/assert"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

type mockGood struct {
	info model.GoodInfo
	qty  int64

	lastUpdateQty int64
}

type mockPortsImpl struct {
	info map[string]mockGood
}

func newMockPortsImpl() *mockPortsImpl {
	return &mockPortsImpl{info: make(map[string]mockGood)}
}

func (m *mockPortsImpl) GetStock(id model.GoodID) model.GoodStock {
	if v, ok := m.info[string(id)]; ok {
		return model.GoodStock{ID: string(id), Quantity: v.qty}
	}
	return model.GoodStock{ID: string(id), Quantity: 0}
}

func (m *mockPortsImpl) GetFreeStock(id model.GoodID) model.GoodStock {
	return m.GetStock(id)
}

func (m *mockPortsImpl) AddGood(id string, name string, description string) {
	m.info[id] = mockGood{
		info: model.GoodInfo{
			ID:          id,
			Name:        name,
			Description: description,
		},
		qty:           0,
		lastUpdateQty: 0,
	}
}

func (m *mockPortsImpl) GetGood(id model.GoodID) *model.GoodInfo {
	if v, ok := m.info[string(id)]; ok {
		return &v.info
	}
	return nil
}

func (m *mockPortsImpl) CreateStockUpdate(ctx context.Context, cmd port.CreateStockUpdateCmd) error {
	for _, v := range cmd.Goods {
		old, ok := m.info[v.Good.ID]
		if !ok {
			return fmt.Errorf("good %s not found", v.Good.ID)
		}
		old.lastUpdateQty = v.Good.Quantity
		old.qty = v.Good.Quantity
		m.info[v.Good.ID] = old
	}

	return nil
}

func TestManageStockService(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	mock := newMockPortsImpl()
	transactionPort := NewMockITransactionPort(ctrl)
	transactionPort.EXPECT().Lock().Times(2)
	transactionPort.EXPECT().Unlock().Times(2)

	app := fx.New(
		fx.Supply(fx.Annotate(mock, fx.As(new(port.ICreateStockUpdatePort)), fx.As(new(port.IGetStockPort)), fx.As(new(port.IGetGoodPort)))),
		fx.Supply(fx.Annotate(transactionPort, fx.As(new(port.ITransactionPort)))),
		fx.Provide(fx.Annotate(NewManageStockService, fx.As(new(port.IAddStockUseCase)), fx.As(new(port.IRemoveStockUseCase)))),
		fx.Invoke(func(lc fx.Lifecycle, addStockUseCase port.IAddStockUseCase, removeStockUseCase port.IRemoveStockUseCase) {
			lc.Append(fx.Hook{OnStart: func(ctx context.Context) error {
				mock.AddGood("1", "hat", "very nice hat")

				addStockCmd := port.AddStockCmd{
					GoodID:   "1",
					Quantity: 10,
				}
				err := addStockUseCase.AddStock(ctx, addStockCmd)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, mock.info["1"].lastUpdateQty, int64(10))
				assert.Equal(t, mock.GetStock("1").Quantity, int64(10))

				remStockCmd := port.RemoveStockCmd{
					GoodID:   "1",
					Quantity: 10,
				}
				err = removeStockUseCase.RemoveStock(ctx, remStockCmd)
				if err != nil {
					t.Error(err)
				}

				assert.Equal(t, mock.GetStock("1").Quantity, int64(0))
				assert.Equal(t, mock.info["1"].lastUpdateQty, int64(0))

				return nil
			}})
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
