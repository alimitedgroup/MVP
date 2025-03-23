package serviceresponse

import (
	"errors"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/stretchr/testify/assert"
)

func TestSetMultipleGoodsQuantityResponse(t *testing.T) {
	obj := NewSetMultipleGoodsQuantityResponse(errors.New("test"), []string{"ciao"})
	assert.Equal(t, "test", obj.GetOperationResult().Error())
	assert.Equal(t, []string{"ciao"}, obj.GetWrongIDSlice())
}

func TestSetGoodQuantityResponse(t *testing.T) {
	obj := NewSetGoodQuantityResponse(errors.New("test"))
	assert.Equal(t, "test", obj.GetOperationResult().Error())
}

func TestGetWarehousesResponse(t *testing.T) {
	obj := NewGetWarehousesResponse(make(map[string]dto.Warehouse))
	assert.Equal(t, make(map[string]dto.Warehouse), obj.GetWarehouseMap())
}

func TestGetGoodsQuantityResponse(t *testing.T) {
	obj := NewGetGoodsQuantityResponse(make(map[string]int64))
	assert.Equal(t, make(map[string]int64), obj.GetMap())
}

func TestGetGoodsInfoResponse(t *testing.T) {
	obj := NewGetGoodsInfoResponse(make(map[string]dto.Good))
	assert.Equal(t, make(map[string]dto.Good), obj.GetMap())
}

func TestAddOrChangeResponse(t *testing.T) {
	obj := NewAddOrChangeResponse(errors.New("test"))
	assert.Equal(t, "test", obj.GetOperationResult().Error())
}
