package persistence

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_stockRepositoryIml.go -package persistence github.com/alimitedgroup/MVP/srv/warehouse/adapter/persistence IStockRepository

func TestStockPersistanceAdapterApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().SetStock(gomock.Any(), gomock.Any()).Return(true).Times(2)

	goods := []model.GoodStock{
		{ID: "1", Quantity: 10},
		{ID: "2", Quantity: 20},
	}

	a := NewStockPersistanceAdapter(mock)

	a.ApplyStockUpdate(goods)
}

func TestStockPersistanceAdapterGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().GetStock(gomock.Any()).Return(int64(10))
	mock.EXPECT().GetFreeStock(gomock.Any()).Return(int64(10))

	a := NewStockPersistanceAdapter(mock)

	require.Equal(t, a.GetStock("1"), model.GoodStock{ID: "1", Quantity: int64(10)})
	require.Equal(t, a.GetFreeStock("1"), model.GoodStock{ID: "1", Quantity: int64(10)})
}

func TestStockPersistanceAdapterGetReserv(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().GetReservation(gomock.Any()).Return(Reservation{Goods: map[string]int64{"1": 10}}, nil)

	a := NewStockPersistanceAdapter(mock)

	reserv, err := a.GetReservation(model.ReservationId("1"))
	require.Nil(t, err)
	require.Equal(t, reserv.ID, model.ReservationId("1"))
	require.Equal(t, reserv.Goods, []model.ReservationGood{{GoodID: "1", Quantity: 10}})
}

func TestStockPersistanceAdapterGetReservNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().GetReservation(gomock.Any()).Return(Reservation{}, ErrReservationNotFound)

	a := NewStockPersistanceAdapter(mock)

	reserv, err := a.GetReservation(model.ReservationId("1"))
	require.NotNil(t, err)
	require.Equal(t, reserv, model.Reservation{})
}

func TestStockPersistanceAdapterApplyReserv(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().ReserveStock(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	a := NewStockPersistanceAdapter(mock)
	require.Nil(t, a.ApplyReservationEvent(model.Reservation{ID: "1", Goods: []model.ReservationGood{{GoodID: "1", Quantity: 10}}}))
}

func TestStockPersistanceAdapterApplyReservErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().ReserveStock(gomock.Any(), gomock.Any(), gomock.Any()).Return(ErrNotEnoughGoods)

	a := NewStockPersistanceAdapter(mock)
	err := a.ApplyReservationEvent(model.Reservation{ID: "1", Goods: []model.ReservationGood{{GoodID: "1", Quantity: 10}}})
	require.NotNil(t, err)
}

func TestStockPersistanceAdapterApplyOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().UnReserveStock(gomock.Any(), gomock.Any()).Return(nil)

	a := NewStockPersistanceAdapter(mock)
	err := a.ApplyOrderFilled(model.Reservation{ID: "1", Goods: []model.ReservationGood{{GoodID: "1", Quantity: 10}}})
	require.Nil(t, err)
}

func TestStockPersistanceAdapterApplyOrderErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIStockRepository(ctrl)

	mock.EXPECT().UnReserveStock(gomock.Any(), gomock.Any()).Return(ErrNotEnoughGoods)

	a := NewStockPersistanceAdapter(mock)
	err := a.ApplyOrderFilled(model.Reservation{ID: "1", Goods: []model.ReservationGood{{GoodID: "1", Quantity: 10}}})
	require.NotNil(t, err)
}
