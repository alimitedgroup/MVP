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

func NewCatalogAdapter(broker *broker.NatsMessageBroker) *CatalogAdapterOut {
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
	//TODO implement me
	panic("implement me")
}

func (c CatalogAdapterOut) ListWarehouses() (map[string]dto.Warehouse, error) {
	//TODO implement me
	panic("implement me")
}

var _ portout.CatalogPortOut = (*CatalogAdapterOut)(nil)
