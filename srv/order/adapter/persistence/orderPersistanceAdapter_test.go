package persistence

import (
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestOrderPersistenceAdapterApplyOrderUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIOrderRepository(ctrl)

	mock.EXPECT().SetOrder(gomock.Any(), gomock.Any()).Return(false)
	mock.EXPECT().GetOrder(gomock.Any()).Return(Order{}, ErrGoodNotFound)

	adapter := NewOrderPersistanceAdapter(mock)

	cmd := port.ApplyOrderUpdateCmd{
		Id:           "1",
		Status:       "Created",
		Name:         "Order 1",
		FullName:     "test test",
		Address:      "via roma 1",
		CreationTime: time.Now().UnixMilli(),
		UpdateTime:   time.Now().UnixMilli(),
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
	err := adapter.ApplyOrderUpdate(cmd)
	require.NoError(t, err)
}

func TestOrderPersistenceAdapterGetOrderExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIOrderRepository(ctrl)

	mock.EXPECT().GetOrder(gomock.Any()).Return(Order{
		ID:           "1",
		Status:       "Created",
		Name:         "Order 1",
		FullName:     "test test",
		Address:      "via roma 1",
		Reservations: []string{},
		Goods: []OrderUpdateGood{
			{
				GoodID:   "1",
				Quantity: 10,
			},
			{
				GoodID:   "2",
				Quantity: 10,
			},
		},
		UpdateTime:   time.Now().UnixMilli(),
		CreationTime: time.Now().UnixMilli(),
	}, nil)

	adapter := NewOrderPersistanceAdapter(mock)

	order, err := adapter.GetOrder("1")
	require.NoError(t, err)
	require.Equal(t, order.Id, model.OrderID("1"))
}

func TestOrderPersistenceAdapterGetOrderNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIOrderRepository(ctrl)

	mock.EXPECT().GetOrder(gomock.Any()).Return(Order{}, ErrOrderNotFound)

	adapter := NewOrderPersistanceAdapter(mock)

	_, err := adapter.GetOrder("1")
	require.Error(t, err, ErrOrderNotFound)
}

func TestOrderPersistenceAdapterGetAllOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIOrderRepository(ctrl)

	mock.EXPECT().GetOrders().Return([]Order{{
		ID:           "1",
		Status:       "Created",
		Name:         "Order 1",
		FullName:     "test test",
		Address:      "via roma 1",
		Reservations: []string{},
		Warehouses:   []OrderWarehouseUsed{},
		Goods: []OrderUpdateGood{
			{
				GoodID:   "1",
				Quantity: 10,
			},
			{
				GoodID:   "2",
				Quantity: 10,
			},
		},
		UpdateTime:   time.Now().UnixMilli(),
		CreationTime: time.Now().UnixMilli(),
	}})

	adapter := NewOrderPersistanceAdapter(mock)

	orders := adapter.GetAllOrder()
	require.Len(t, orders, 1)
	require.Equal(t, orders[0].Id, model.OrderID("1"))
}

func TestOrderPersistenceAdapterSetComplete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIOrderRepository(ctrl)

	mock.EXPECT().SetComplete(gomock.Any()).Return(nil)

	adapter := NewOrderPersistanceAdapter(mock)

	err := adapter.SetComplete("1")
	require.NoError(t, err)
}

func TestOrderPersistenceAdapterSetCompleteErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIOrderRepository(ctrl)

	mock.EXPECT().SetComplete(gomock.Any()).Return(ErrOrderNotFound)

	adapter := NewOrderPersistanceAdapter(mock)

	err := adapter.SetComplete("1")
	require.Error(t, err, ErrOrderNotFound)
}

func TestOrderPersistenceAdapterAddCompletedWarehouse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIOrderRepository(ctrl)

	mock.EXPECT().AddCompletedWarehouse(gomock.Any(), gomock.Any(), gomock.Any()).Return(Order{
		ID:           "1",
		Status:       "Created",
		Name:         "Order 1",
		FullName:     "test test",
		Address:      "via roma 1",
		Reservations: []string{},
		Warehouses:   []OrderWarehouseUsed{},
		Goods: []OrderUpdateGood{
			{
				GoodID:   "1",
				Quantity: 10,
			},
			{
				GoodID:   "2",
				Quantity: 10,
			},
		},
		UpdateTime:   time.Now().UnixMilli(),
		CreationTime: time.Now().UnixMilli(),
	}, nil)

	adapter := NewOrderPersistanceAdapter(mock)

	cmd := port.SetCompletedWarehouseCmd{
		OrderId:     "1",
		WarehouseId: "1",
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
	order, err := adapter.SetCompletedWarehouse(cmd)
	require.NoError(t, err)
	require.Equal(t, order.Id, model.OrderID("1"))

}

func TestOrderPersistenceAdapterAddCompletedWarehouseErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIOrderRepository(ctrl)

	mock.EXPECT().AddCompletedWarehouse(gomock.Any(), gomock.Any(), gomock.Any()).Return(Order{}, ErrGoodNotFound)

	adapter := NewOrderPersistanceAdapter(mock)

	cmd := port.SetCompletedWarehouseCmd{
		OrderId:     "1",
		WarehouseId: "1",
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
	order, err := adapter.SetCompletedWarehouse(cmd)
	require.Error(t, err, ErrGoodNotFound)
	require.Equal(t, order.Id, model.OrderID(""))
}
