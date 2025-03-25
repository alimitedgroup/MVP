package servicecmd

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetMultipleGoodsQuantityCmd(t *testing.T) {
	obj := NewSetMultipleGoodsQuantityCmd("test", []stream.StockUpdateGood{{GoodID: "test", Quantity: 7, Delta: 0}})
	assert.Equal(t, []stream.StockUpdateGood{{GoodID: "test", Quantity: 7, Delta: 0}}, obj.GetGoods())
	assert.Equal(t, "test", obj.GetWarehouseID())
}

func TestAddChangeGoodCmd(t *testing.T) {
	obj := NewAddChangeGoodCmd("test-id", "test-name", "test-description")
	assert.Equal(t, "test-id", obj.GetId())
	assert.Equal(t, "test-name", obj.GetName())
	assert.Equal(t, "test-description", obj.GetDescription())
}

func TestSetGoodQuantityCmd(t *testing.T) {
	obj := NewSetGoodQuantityCmd("test", "test-good", 7)
	assert.Equal(t, "test", obj.GetWarehouseId())
	assert.Equal(t, "test-good", obj.GetGoodId())
	assert.Equal(t, int64(7), obj.GetNewQuantity())
}
func TestGetWarehousesCmd(t *testing.T) {
	obj := NewGetWarehousesCmd()
	require.NotNil(t, obj)
}
func TestGetGoodsQuantityCmd(t *testing.T) {
	obj := NewGetGoodsQuantityCmd()
	require.NotNil(t, obj)
}
func TestGetGoodsInfoCmd(t *testing.T) {
	obj := NewGetGoodsInfoCmd()
	require.NotNil(t, obj)
}
