package persistence

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestStockPersistenceAdapterApplyStockUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().SetStock(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).Times(2)

	adapter := NewStockPersistanceAdapter(mock)

	cmd := port.ApplyStockUpdateCmd{
		WarehouseID: "1",
		Goods: []model.GoodStock{
			{
				ID:       "1",
				Quantity: 10,
			},
			{
				ID:       "2",
				Quantity: 10,
			},
		},
	}
	err := adapter.ApplyStockUpdate(cmd)
	require.NoError(t, err)
}

func TestStockPersistenceAdapterGetWarehouses(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().GetWarehouses().Return([]string{"1", "2"})

	adapter := NewStockPersistanceAdapter(mock)

	warehoues := adapter.GetWarehouses()
	require.Len(t, warehoues, 2)
}

func TestStockPersistenceAdapterGetGlobalStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().GetGlobalStock(gomock.Any()).Return(int64(10))

	adapter := NewStockPersistanceAdapter(mock)

	stock := adapter.GetGlobalStock(model.GoodID("1"))
	require.Equal(t, stock.Quantity, int64(10))
}

func TestStockPersistenceAdapterGetStockExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().GetStock(gomock.Any(), gomock.Any()).Return(int64(10), nil)

	adapter := NewStockPersistanceAdapter(mock)

	stock, err := adapter.GetStock(port.GetStockCmd{
		WarehouseID: "1",
		GoodID:      "1",
	})
	require.NoError(t, err)
	require.Equal(t, stock.Quantity, int64(10))
}

func TestStockPersistenceAdapterGetStockNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().GetStock(gomock.Any(), gomock.Any()).Return(int64(0), ErrGoodNotFound)

	adapter := NewStockPersistanceAdapter(mock)

	stock, err := adapter.GetStock(port.GetStockCmd{
		WarehouseID: "1",
		GoodID:      "1",
	})
	require.Error(t, err, ErrGoodNotFound)
	require.Equal(t, stock.Quantity, int64(0))
}

func TestStockPersistenceAdapterGetStockWarehouseNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().GetStock(gomock.Any(), gomock.Any()).Return(int64(0), ErrWarehouseNotFound)

	adapter := NewStockPersistanceAdapter(mock)

	stock, err := adapter.GetStock(port.GetStockCmd{
		WarehouseID: "1",
		GoodID:      "1",
	})
	require.Error(t, err, ErrWarehouseNotFound)
	require.Equal(t, stock.Quantity, int64(0))
}
