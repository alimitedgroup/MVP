package adapterout

import (
	"encoding/json"
	"fmt"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/nats-io/nats.go"
)

type OrderAdapterOut struct {
	Broker *broker.NatsMessageBroker
}

func NewOrderAdapter(broker *broker.NatsMessageBroker) portout.OrderPortOut {
	return &OrderAdapterOut{
		Broker: broker,
	}
}

func (c OrderAdapterOut) CreateOrder(dto request.CreateOrderRequestDTO) (response.OrderCreateInfo, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return response.OrderCreateInfo{}, err
	}

	resp, err := c.Broker.Nats.Request("order.create", payload, nats.DefaultTimeout)
	if err != nil {
		return response.OrderCreateInfo{}, err
	}

	var respDto response.OrderCreateResponseDTO
	err = json.Unmarshal(resp.Data, &respDto)
	if err != nil {
		return response.OrderCreateInfo{}, err
	}

	if respDto.Error != "" {
		return response.OrderCreateInfo{}, fmt.Errorf("%s", respDto.Error)
	}

	return respDto.Message, err
}

func (c OrderAdapterOut) GetAllOrders() ([]response.OrderInfo, error) {
	resp, err := c.Broker.Nats.Request("order.get.all", []byte("{}"), nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var respDto response.GetAllOrderResponseDTO
	err = json.Unmarshal(resp.Data, &respDto)
	if err != nil {
		return nil, err
	}

	if respDto.Error != "" {
		return nil, fmt.Errorf("%s", respDto.Error)
	}

	return respDto.Message, err
}

func (c OrderAdapterOut) CreateTransfer(dto request.CreateTransferRequestDTO) (response.TransferCreateInfo, error) {
	payload, err := json.Marshal(dto)
	if err != nil {
		return response.TransferCreateInfo{}, err
	}

	resp, err := c.Broker.Nats.Request("transfer.create", payload, nats.DefaultTimeout)
	if err != nil {
		return response.TransferCreateInfo{}, err
	}

	var respDto response.TransferCreateResponseDTO
	err = json.Unmarshal(resp.Data, &respDto)
	if err != nil {
		return response.TransferCreateInfo{}, err
	}

	if respDto.Error != "" {
		return response.TransferCreateInfo{}, fmt.Errorf("%s", respDto.Error)
	}

	return respDto.Message, err
}

func (c OrderAdapterOut) GetAllTransfers() ([]response.TransferInfo, error) {
	resp, err := c.Broker.Nats.Request("transfer.get.all", []byte("{}"), nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var respDto response.GetAllTransferResponseDTO
	err = json.Unmarshal(resp.Data, &respDto)
	if err != nil {
		return nil, err
	}

	if respDto.Error != "" {
		return nil, fmt.Errorf("%s", respDto.Error)
	}

	return respDto.Message, err
}

var _ portout.OrderPortOut = (*OrderAdapterOut)(nil)
