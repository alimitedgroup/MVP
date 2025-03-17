package adapterout

import (
	"encoding/json"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
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
	resp, err := c.Broker.Nats.Request("catalog.getGoodsGlobalQuantity", []byte("{}"), nats.DefaultTimeout)
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

var _ portout.CatalogPortOut = (*CatalogAdapterOut)(nil)
