package adapterin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/stretchr/testify/require"
)

func TestGetWarehouses(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().GetWarehouses().Return(
		[]types.WarehouseOverview{{ID: "id1"}, {ID: "id2"}, {ID: "id3"}},
		nil,
	)
	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)

	req, err := http.NewRequest(http.MethodGet, s.base+"/api/v1/warehouses", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")
	resp, err := client.Do(req)
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
	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)

	req, err := http.NewRequest(http.MethodGet, s.base+"/api/v1/warehouses", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var respbody response.ResponseDTO[any]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "internal_error", respbody.Error)
}
