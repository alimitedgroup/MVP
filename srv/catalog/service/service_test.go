package service

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/catalog/catalogCommon"
	service_Cmd "github.com/alimitedgroup/MVP/srv/catalog/service/Cmd"
	service_Response "github.com/alimitedgroup/MVP/srv/catalog/service/Response"
	service_portOut "github.com/alimitedgroup/MVP/srv/catalog/service/portOut"
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

func (fp *fakeAddOrChangeGoodDataPort) AddOrChangeGoodData(agc *service_Cmd.AddChangeGoodCmd) *service_Response.AddOrChangeResponse {
	if agc.GetId() == "test-wrong-ID" {
		return service_Response.NewAddOrChangeResponse("Not a valid goodID")
	}
	return service_Response.NewAddOrChangeResponse("Success")
}

func (fp *fakeSetGoodQuantityPort) SetGoodQuantity(agqc *service_Cmd.SetGoodQuantityCmd) *service_Response.SetGoodQuantityResponse {
	if agqc.GetGoodId() == "test-wrong-ID" {
		return service_Response.NewSetGoodQuantityResponse("Not a valid goodID")
	}
	return service_Response.NewSetGoodQuantityResponse("Success")
}

func (fp *fakeGetGoodsQuantityPort) GetGoodsQuantity(ggqc *service_Cmd.GetGoodsQuantityCmd) *service_Response.GetGoodsQuantityResponse {
	goodMap := map[string]int64{}
	goodMap["test-ID"] = int64(7)
	return service_Response.NewGetGoodsQuantityResponse(goodMap)
}

func (fp *fakeGetGoodsInfoPort) GetGoodsInfo(ggqc *service_Cmd.GetGoodsInfoCmd) *service_Response.GetGoodsInfoResponse {
	goods := make(map[string]catalogCommon.Good)
	goods["test-ID"] = *catalogCommon.NewGood("test-ID", "test-name", "test-description")
	return service_Response.NewGetGoodsInfoResponse(goods)
}

func (fp *fakeGetWarehousesPort) GetWarehouses(gwc *service_Cmd.GetWarehousesCmd) *service_Response.GetWarehousesResponse {
	warehouses := make(map[string]catalogCommon.Warehouse)
	warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")
	return service_Response.NewGetWarehousesResponse(warehouses)
}

// FINE DESCRIZIONE PORTE MOCK

func TestGetWarehouses(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(service_portOut.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(service_portOut.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(service_portOut.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(service_portOut.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(service_portOut.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			warehouses := make(map[string]catalogCommon.Warehouse)
			warehouses["test-warehouse-ID"] = *catalogCommon.NewWarehouse("test-warehose-ID")
			cmd := service_Cmd.NewGetWarehousesCmd()
			response := cs.GetWarehouses(cmd)
			assert.Equal(t, response.GetWarehouseMap(), warehouses)
		}),
	)
}

func TestAddOrChangeGoodData(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(service_portOut.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(service_portOut.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(service_portOut.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(service_portOut.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(service_portOut.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			response := cs.AddOrChangeGoodData(service_Cmd.NewAddChangeGoodCmd("test-ID", "test-name", "test-description"))
			assert.Equal(t, response.GetOperationResult(), "Success")
		}),
	)
}

func TestAddOrChangeGoodData_WrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(service_portOut.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(service_portOut.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(service_portOut.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(service_portOut.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(service_portOut.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			response := cs.AddOrChangeGoodData(service_Cmd.NewAddChangeGoodCmd("test-wrong-ID", "test-name", "test-description"))
			assert.Equal(t, response.GetOperationResult(), "Not a valid goodID")
		}),
	)
}

func TestSetMultipleGoodsQuantity(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(service_portOut.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(service_portOut.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(service_portOut.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(service_portOut.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(service_portOut.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			goods := []stream.StockUpdateGood{}
			goods = append(goods, stream.StockUpdateGood{GoodID: "test-ID", Quantity: int64(7), Delta: int64(0)})
			goods = append(goods, stream.StockUpdateGood{GoodID: "2test-ID", Quantity: int64(9), Delta: int64(1)})
			cmd := service_Cmd.NewMultipleGoodsQuantityCmd("test-warehouse-ID", goods)
			response := cs.SetMultipleGoodsQuantity(cmd)
			assert.Equal(t, response.GetOperationResult(), "Success")
			assert.Equal(t, len(response.GetWrongIDSlice()), 0)
		}),
	)
}

func TestSetMultipleGoodsQuantityWithWrongID(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(service_portOut.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(service_portOut.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(service_portOut.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(service_portOut.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(service_portOut.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			goods := []stream.StockUpdateGood{}
			goods = append(goods, stream.StockUpdateGood{GoodID: "test-wrong-ID", Quantity: int64(7), Delta: int64(0)})
			goods = append(goods, stream.StockUpdateGood{GoodID: "2test-ID", Quantity: int64(9), Delta: int64(1)})
			cmd := service_Cmd.NewMultipleGoodsQuantityCmd("test-warehouse-ID", goods)
			response := cs.SetMultipleGoodsQuantity(cmd)
			assert.Equal(t, response.GetOperationResult(), "Errors")
			assert.Equal(t, len(response.GetWrongIDSlice()), 1)
			assert.Equal(t, response.GetWrongIDSlice()[0], "test-wrong-ID")
		}),
	)
}

func TestGetGoodsQuantity(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(service_portOut.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(service_portOut.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(service_portOut.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(service_portOut.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(service_portOut.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			response := cs.GetGoodsQuantity(service_Cmd.NewGetGoodsQuantityCmd())
			assert.Equal(t, response.GetMap()["test-ID"], int64(7))
		}),
	)
}

func TestGetGoodsInfo(t *testing.T) {
	fx.New(
		fx.Provide(
			fx.Annotate(NewFakeAddOrChangeGoodDataPort,
				fx.As(new(service_portOut.IAddOrChangeGoodDataPort)),
			),
			fx.Annotate(NewFakeGetGoodsInfoPort,
				fx.As(new(service_portOut.IGetGoodsInfoPort)),
			),
			fx.Annotate(NewFakeGetGoodsQuantityPort,
				fx.As(new(service_portOut.IGetGoodsQuantityPort)),
			),
			fx.Annotate(
				NewFakeSetGoodQuantityPort,
				fx.As(new(service_portOut.ISetGoodQuantityPort)),
			),
			fx.Annotate(
				NewFakeGetWarehousesPort,
				fx.As(new(service_portOut.IGetWarehousesInfoPort)),
			),
		),
		fx.Provide(NewCatalogService),
		fx.Invoke(func(cs *CatalogService) {
			response := cs.GetGoodsInfo(service_Cmd.NewGetGoodsInfoCmd())
			assert.Equal(t, response.GetMap()["test-ID"].GetID(), "test-ID")
			assert.Equal(t, response.GetMap()["test-ID"].GetName(), "test-name")
			assert.Equal(t, response.GetMap()["test-ID"].GetDescription(), "test-description")
		}),
	)
}
