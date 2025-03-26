package adapterin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateOrder(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.order.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)

	payload := bytes.NewReader([]byte(`{"name":"Order 1","full_name":"John Doe","address":"123 Main St","goods":{"1":1}}`))
	resp, err := client.Post(s.base+"/api/v1/orders", "application/json", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.CreateOrderResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.NotEmpty(t, respbody.OrderID)
}
