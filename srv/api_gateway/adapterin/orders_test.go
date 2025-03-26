package adapterin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/stretchr/testify/require"
)

func TestGetOrders(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.order.EXPECT().GetOrders().Return(
		[]dto.Order{
			{
				OrderID:  "1",
				Status:   "Created",
				Name:     "Order 1",
				FullName: "John Doe",
				Address:  "123 Main St",
				Goods:    map[string]int64{"1": 1},
			},
		},
		nil,
	)

	resp, err := client.Get(s.base + "/api/v1/orders")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.GetOrdersResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, []dto.Order{
		{
			OrderID:  "1",
			Status:   "Created",
			Name:     "Order 1",
			FullName: "John Doe",
			Address:  "123 Main St",
			Goods:    map[string]int64{"1": 1},
		},
	}, respbody.Orders)
}

func TestGetOrdersError(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.order.EXPECT().GetOrders().Return([]dto.Order{}, fmt.Errorf("some error"))

	resp, err := client.Get(s.base + "/api/v1/orders")
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "internal_error", respbody.Error)
}
