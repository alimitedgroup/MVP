package business

import (
	"fmt"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestGetWarehouses(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	catalog.EXPECT().ListWarehouses().Return(map[string]dto.Warehouse{
		"abc": {ID: "abc", Stock: map[string]int64{"id1": 20}},
		"def": {ID: "def", Stock: map[string]int64{"id1": 10, "id2": 20}},
	}, nil)

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	warehouses, err := business.GetWarehouses()
	require.NoError(t, err)
	require.Len(t, warehouses, 2)
	require.ElementsMatch(t, []portin.WarehouseOverview{{ID: "abc"}, {ID: "def"}}, warehouses)
}

func TestGetWarehousesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	catalog.EXPECT().ListWarehouses().Return(nil, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	warehouses, err := business.GetWarehouses()
	require.Nil(t, warehouses)
	require.ErrorIs(t, err, ErrorGetWarehouses)
}

func TestGetGoods(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	catalog.EXPECT().ListGoods().Return(map[string]dto.Good{
		"id1": {Name: "abc", Description: "abcdesc", ID: "id1"},
		"id2": {Name: "def", Description: "defdesc", ID: "id2"},
	}, nil)
	catalog.EXPECT().ListStock().Return(
		map[string]int64{"id1": 20, "id2": 10},
		nil,
	)
	catalog.EXPECT().ListWarehouses().Return(
		map[string]dto.Warehouse{
			"warehouse1": {ID: "warehouse1", Stock: map[string]int64{"id1": 20, "id2": 10}},
		},
		nil,
	)

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	goods, err := business.GetGoods()
	require.NoError(t, err)
	require.Len(t, goods, 2)
	require.ElementsMatch(t, []dto.GoodAndAmount{
		{ID: "id1", Amount: 20, Name: "abc", Description: "abcdesc", Amounts: map[string]int64{"warehouse1": 20}},
		{ID: "id2", Amount: 10, Name: "def", Description: "defdesc", Amounts: map[string]int64{"warehouse1": 10}},
	}, goods)
}

func TestGetGoodsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	catalog.EXPECT().ListGoods().Return(nil, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	goods, err := business.GetGoods()
	require.Nil(t, goods)
	require.ErrorIs(t, err, ErrorGetGoods)
}

func TestGetGoodsStockError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	catalog.EXPECT().ListGoods().Return(map[string]dto.Good{
		"id1": {Name: "abc", Description: "abcdesc", ID: "id1"},
		"id2": {Name: "def", Description: "defdesc", ID: "id2"},
	}, nil)
	catalog.EXPECT().ListStock().Return(nil, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	goods, err := business.GetGoods()
	require.Nil(t, goods)
	require.ErrorIs(t, err, ErrorGetStock)
}

func TestGetGoodsMissingStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	catalog.EXPECT().ListGoods().Return(map[string]dto.Good{
		"id1": {Name: "abc", Description: "abcdesc", ID: "id1"},
		"id2": {Name: "def", Description: "defdesc", ID: "id2"},
	}, nil)
	catalog.EXPECT().ListStock().Return(
		map[string]int64{"id1": 20},
		nil,
	)
	catalog.EXPECT().ListWarehouses().Return(map[string]dto.Warehouse{
		"warehouse1": {ID: "warehouse1", Stock: map[string]int64{"id1": 15}},
		"warehouse2": {ID: "warehouse2", Stock: map[string]int64{"id1": 5}},
	}, nil)

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	goods, err := business.GetGoods()
	require.ElementsMatch(t, goods, []dto.GoodAndAmount{
		{Name: "abc", Description: "abcdesc", ID: "id1", Amount: 20, Amounts: map[string]int64{"warehouse1": 15, "warehouse2": 5}},
		{Name: "def", Description: "defdesc", ID: "id2", Amount: 0, Amounts: map[string]int64{}},
	})
	require.NoError(t, err)
}
