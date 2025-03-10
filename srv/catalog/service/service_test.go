package service

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	servicecmd "github.com/alimitedgroup/MVP/srv/catalog/service/cmd"
	serviceportout "github.com/alimitedgroup/MVP/srv/catalog/service/portout"
	serviceresponse "github.com/alimitedgroup/MVP/srv/catalog/service/response"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

// INIZIO DESCRIZIONE PORTE MOCK

type fakeAddOrChangeGoodDataPort struct {
}

type fakeSetGoodQuantityPort struct {
}

type fakeGetGoodsQuantityPort struct {
}

type fakeGetGoodsInfoPort struct {
}

type fakeGetWarehousesPort struct {
}

func NewFakeAddOrChangeGoodDataPort() *fakeAddOrChangeGoodDataPort {
	return &fakeAddOrChangeGoodDataPort{}
}

func NewFakeSetGoodQuantityPort() *fakeSetGoodQuantityPort {
	return &fakeSetGoodQuantityPort{}
}

func NewFakeGetGoodsQuantityPort() *fakeGetGoodsQuantityPort {
	return &fakeGetGoodsQuantityPort{}
}

func NewFakeGetGoodsInfoPort() *fakeGetGoodsInfoPort {
	return &fakeGetGoodsInfoPort{}
}

func NewFakeGetWarehousesPort() *fakeGetWarehousesPort {
	return &fakeGetWarehousesPort{}
}

func (fp *fakeAddOrChangeGoodDataPort) AddOrChangeGoodData(agc *servicecmd.AddChangeGoodCmd) *serviceresponse.AddOrChangeResponse {
	if agc.GetId() == "test-wrong-ID" {
		return serviceresponse.NewAddOrChangeResponse(catalogCommon.ErrGoodIdNotValid)
	}
	return serviceresponse.NewAddOrChangeResponse(nil)
}

func (fp *fakeSetGoodQuantityPort) SetGoodQuantity(agqc *servicecmd.SetGoodQuantityCmd) *serviceresponse.SetGoodQuantityResponse {
	if agqc.GetGoodId() == "test-wrong-ID" {
		return serviceresponse.NewSetGoodQuantityResponse(catalogCommon.ErrGoodIdNotValid)
	}
	return serviceresponse.NewSetGoodQuantityResponse(nil)
}

func (fp *fakeGetGoodsQuantityPort) GetGoodsQuantity(ggqc *servicecmd.GetGoodsQuantityCmd) *serviceresponse.GetGoodsQuantityResponse {
	goodMap := map[string]int64{}
	goodMap["test-ID"] = int64(7)
	return serviceresponse.NewGetGoodsQuantityResponse(goodMap)
}

func (fp *fakeGetGoodsInfoPort) GetGoodsInfo(ggqc *servicecmd.GetGoodsInfoCmd) *serviceresponse.GetGoodsInfoResponse {
	goods := make(map[string]catalogCommon.Good)
	goods["test-ID"] = *catalogCommon.NewGood("test-ID", "test-name", "test-description")
	return serviceresponse.NewGetGoodsInfoResponse(goods)
}

func (fp *fakeGetWarehousesPort) GetWarehouses(gwc *servicecmd.GetWarehousesCmd) *serviceresponse.GetWarehousesResponse {
	warehouses := make(map[string]catalogCommon.Warehouse)
	warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")
	return serviceresponse.NewGetWarehousesResponse(warehouses)
}

// FINE DESCRIZIONE PORTE MOCK

func TestGetWarehouses(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(serviceportout.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(serviceportout.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(serviceportout.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(serviceportout.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(serviceportout.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			warehouses := make(map[string]catalogCommon.Warehouse)
			warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")
			cmd := servicecmd.NewGetWarehousesCmd()
			response := cs.GetWarehouses(cmd)
			assert.Equal(t, response.GetWarehouseMap(), warehouses)
		}),
	)
}

func TestAddOrChangeGoodData(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(serviceportout.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(serviceportout.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(serviceportout.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(serviceportout.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(serviceportout.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			response := cs.AddOrChangeGoodData(servicecmd.NewAddChangeGoodCmd("test-ID", "test-name", "test-description"))
			assert.Equal(t, response.GetOperationResult(), nil)
		}),
	)
}

func TestAddOrChangeGoodData_WrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(serviceportout.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(serviceportout.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(serviceportout.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(serviceportout.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(serviceportout.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			response := cs.AddOrChangeGoodData(servicecmd.NewAddChangeGoodCmd("test-wrong-ID", "test-name", "test-description"))
			assert.Equal(t, response.GetOperationResult(), catalogCommon.ErrGoodIdNotValid)
		}),
	)
}

func TestSetMultipleGoodsQuantity(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(serviceportout.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(serviceportout.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(serviceportout.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(serviceportout.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(serviceportout.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			goods := []stream.StockUpdateGood{}
			goods = append(goods, stream.StockUpdateGood{GoodID: "test-ID", Quantity: int64(7), Delta: int64(0)})
			goods = append(goods, stream.StockUpdateGood{GoodID: "2test-ID", Quantity: int64(9), Delta: int64(1)})
			cmd := servicecmd.NewSetMultipleGoodsQuantityCmd("test-warehouse-ID", goods)
			response := cs.SetMultipleGoodsQuantity(cmd)
			assert.Equal(t, response.GetOperationResult(), nil)
			assert.Equal(t, len(response.GetWrongIDSlice()), 0)
		}),
	)
}

func TestSetMultipleGoodsQuantityWithWrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(serviceportout.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(serviceportout.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(serviceportout.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(serviceportout.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(serviceportout.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			goods := []stream.StockUpdateGood{}
			goods = append(goods, stream.StockUpdateGood{GoodID: "test-wrong-ID", Quantity: int64(7), Delta: int64(0)})
			goods = append(goods, stream.StockUpdateGood{GoodID: "2test-ID", Quantity: int64(9), Delta: int64(1)})
			cmd := servicecmd.NewSetMultipleGoodsQuantityCmd("test-warehouse-ID", goods)
			response := cs.SetMultipleGoodsQuantity(cmd)
			assert.Equal(t, response.GetOperationResult(), catalogCommon.ErrGenericFailure)
			assert.Equal(t, len(response.GetWrongIDSlice()), 1)
			assert.Equal(t, response.GetWrongIDSlice()[0], "test-wrong-ID")
		}),
	)
}

func TestGetGoodsQuantity(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(serviceportout.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(serviceportout.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(serviceportout.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(serviceportout.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(serviceportout.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			response := cs.GetGoodsQuantity(servicecmd.NewGetGoodsQuantityCmd())
			assert.Equal(t, response.GetMap()["test-ID"], int64(7))
		}),
	)
}

func TestGetGoodsInfo(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(serviceportout.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(serviceportout.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(serviceportout.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(serviceportout.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(serviceportout.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			response := cs.GetGoodsInfo(servicecmd.NewGetGoodsInfoCmd())
			result := response.GetMap()["test-ID"]
			assert.Equal(t, result.GetID(), "test-ID")
			assert.Equal(t, result.GetName(), "test-name")
			assert.Equal(t, result.GetDescription(), "test-description")
		}),
	)
}
