package adapterout

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/common/stream"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type CatalogAdapterOut struct {
	Broker *broker.NatsMessageBroker
}

func NewCatalogAdapter(broker *broker.NatsMessageBroker) portout.CatalogPortOut {
	return &CatalogAdapterOut{
		Broker: broker,
	}
}

func (c CatalogAdapterOut) ListGoods() (map[string]dto.Good, error) {
	resp, err := c.Broker.Nats.Request("catalog.getGoods", []byte("{}"), nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var goods dto.GetGoodsDataResponseDTO
	err = json.Unmarshal(resp.Data, &goods)
	if err != nil {
		return nil, err
	}

	if goods.Err != "" {
		return nil, fmt.Errorf("%s", goods.Err)
	}

	return goods.GoodMap, err
}

func (c CatalogAdapterOut) ListStock() (map[string]int64, error) {
	resp, err := c.Broker.Nats.Request("catalog.getGoodsGlobalQuantity", []byte("{}"), nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var goods dto.GetGoodsQuantityResponseDTO
	err = json.Unmarshal(resp.Data, &goods)
	if err != nil {
		return nil, err
	}

	if goods.Err != "" {
		return nil, fmt.Errorf("%s", goods.Err)
	}

	return goods.GoodMap, err
}

func (c CatalogAdapterOut) ListWarehouses() (map[string]dto.Warehouse, error) {
	resp, err := c.Broker.Nats.Request("catalog.getWarehouses", []byte("{}"), nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var goods dto.GetWarehouseResponseDTO
	err = json.Unmarshal(resp.Data, &goods)
	if err != nil {
		return nil, err
	}

	if goods.Err != "" {
		return nil, fmt.Errorf("%s", goods.Err)
	}

	return goods.WarehouseMap, err
}

func (c CatalogAdapterOut) AddStock(warehouseId string, goodId string, quantity int64) error {
	payload, err := json.Marshal(request.AddStockRequestDTO{
		GoodID:   goodId,
		Quantity: quantity,
	})
	if err != nil {
		return err
	}

	resp, err := c.Broker.Nats.Request(fmt.Sprintf("warehouse.%s.stock.add", warehouseId), payload, nats.DefaultTimeout)
	if err != nil {
		return err
	}

	var respDto response.ResponseDTO[string]
	err = json.Unmarshal(resp.Data, &respDto)
	if err != nil {
		return err
	}

	if respDto.Error != "" {
		return fmt.Errorf("%s", respDto.Error)
	}

	return err
}

func (c CatalogAdapterOut) RemoveStock(warehouseId string, goodId string, quantity int64) error {
	payload, err := json.Marshal(request.RemoveStockRequestDTO{
		GoodID:   goodId,
		Quantity: quantity,
	})
	if err != nil {
		return err
	}

	resp, err := c.Broker.Nats.Request(fmt.Sprintf("warehouse.%s.stock.remove", warehouseId), payload, nats.DefaultTimeout)
	if err != nil {
		return err
	}

	var respDto response.ResponseDTO[string]
	err = json.Unmarshal(resp.Data, &respDto)
	if err != nil {
		return err
	}

	if respDto.Error != "" {
		return fmt.Errorf("%s", respDto.Error)
	}

	return err
}

func (c CatalogAdapterOut) CreateGood(ctx context.Context, name string, description string) (string, error) {
	goodId := uuid.New().String()

	err := c.UpdateGood(ctx, goodId, name, description)
	if err != nil {
		return "", err
	}

	return goodId, err
}

func (c CatalogAdapterOut) UpdateGood(ctx context.Context, goodId string, name string, description string) error {
	req := stream.GoodUpdateData{
		GoodID:             goodId,
		GoodNewName:        name,
		GoodNewDescription: description,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.Broker.Js.Publish(ctx, "good.update", reqBytes)
	if err != nil {
		return err
	}

	_ = resp

	return err
}

var _ portout.CatalogPortOut = (*CatalogAdapterOut)(nil)
