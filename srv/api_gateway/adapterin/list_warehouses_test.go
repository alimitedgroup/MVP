package adapterin

import (
	"encoding/json"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetWarehouses(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().GetWarehouses().Return(
		[]portin.WarehouseOverview{{"id1"}, {"id2"}, {"id3"}},
		nil,
	)

	resp, err := client.Get(s.base + "/api/v1/warehouses")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.GetWarehousesResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, []string{"id1", "id2", "id3"}, respbody.Ids)
}

func TestGetWarehousesError(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().GetWarehouses().Return(
		nil,
		fmt.Errorf("some error"),
	)

	resp, err := client.Get(s.base + "/api/v1/warehouses")
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var respbody response.ResponseDTO[any]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "internal_error", respbody.Error)
}
