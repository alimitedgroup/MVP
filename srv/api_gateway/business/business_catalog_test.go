package business

import (
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetWarehouses(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)

	catalog.EXPECT().ListWarehouses().Return(map[string]dto.Warehouse{
		"abc": {ID: "abc", Stock: map[string]int64{"id1": 20}},
		"def": {ID: "def", Stock: map[string]int64{"id1": 10, "id2": 20}},
	}, nil)

	business := NewBusiness(auth, catalog)
	warehouses, err := business.GetWarehouses()
	require.NoError(t, err)
	require.Len(t, warehouses, 2)
	require.Equal(t, "abc", warehouses[0].ID)
	require.Equal(t, "def", warehouses[1].ID)
}

func TestGetWarehousesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)

	catalog.EXPECT().ListWarehouses().Return(nil, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog)
	warehouses, err := business.GetWarehouses()
	require.Nil(t, warehouses)
	require.ErrorIs(t, err, ErrorGetWarehouses)
}

func TestGetGoods(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)

	catalog.EXPECT().ListGoods().Return(map[string]dto.Good{
		"id1": {Name: "abc", Description: "abcdesc", ID: "id1"},
		"id2": {Name: "def", Description: "defdesc", ID: "id2"},
	}, nil)
	catalog.EXPECT().ListStock().Return(
		map[string]int64{"id1": 20, "id2": 10},
		nil,
	)

	business := NewBusiness(auth, catalog)
	goods, err := business.GetGoods()
	require.NoError(t, err)
	require.Len(t, goods, 2)

	var id1, id2 int
	if goods[0].ID == "id1" {
		id1 = 0
		id2 = 1
	} else {
		id1 = 1
		id2 = 0
	}
	require.Equal(t, "abc", goods[id1].Name)
	require.Equal(t, "id1", goods[id1].ID)
	require.Equal(t, "abcdesc", goods[id1].Description)
	require.Equal(t, int64(20), goods[id1].Amount)
	require.Equal(t, "def", goods[id2].Name)
	require.Equal(t, "id2", goods[id2].ID)
	require.Equal(t, "defdesc", goods[id2].Description)
	require.Equal(t, int64(10), goods[id2].Amount)
}

func TestGetGoodsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)

	catalog.EXPECT().ListGoods().Return(nil, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog)
	goods, err := business.GetGoods()
	require.Nil(t, goods)
	require.ErrorIs(t, err, ErrorGetGoods)
}

func TestGetGoodsStockError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)

	catalog.EXPECT().ListGoods().Return(map[string]dto.Good{
		"id1": {Name: "abc", Description: "abcdesc", ID: "id1"},
		"id2": {Name: "def", Description: "defdesc", ID: "id2"},
	}, nil)
	catalog.EXPECT().ListStock().Return(nil, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog)
	goods, err := business.GetGoods()
	require.Nil(t, goods)
	require.ErrorIs(t, err, ErrorGetStock)
}

func TestGetGoodsMissingStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)

	catalog.EXPECT().ListGoods().Return(map[string]dto.Good{
		"id1": {Name: "abc", Description: "abcdesc", ID: "id1"},
		"id2": {Name: "def", Description: "defdesc", ID: "id2"},
	}, nil)
	catalog.EXPECT().ListStock().Return(
		map[string]int64{"id1": 20},
		nil,
	)

	business := NewBusiness(auth, catalog)
	goods, err := business.GetGoods()
	require.Nil(t, goods)
	require.ErrorIs(t, err, ErrorGoodWithoutStock)
}
