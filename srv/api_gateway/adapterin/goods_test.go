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

func TestGetGoods(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().GetGoods().Return(
		[]dto.GoodAndAmount{
			{Name: "Apple", Description: "A tasty apple", ID: "id1", Amount: 20, Amounts: map[string]int64{"id1": 20}},
			{Name: "Orange", Description: "A tasty orange", ID: "id2", Amount: 10, Amounts: map[string]int64{"id1": 10}},
		},
		nil,
	)
	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)

	req, err := http.NewRequest(http.MethodGet, s.base+"/api/v1/goods", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.GetGoodsResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, []dto.GoodAndAmount{
		{Name: "Apple", Description: "A tasty apple", ID: "id1", Amount: 20, Amounts: map[string]int64{"id1": 20}},
		{Name: "Orange", Description: "A tasty orange", ID: "id2", Amount: 10, Amounts: map[string]int64{"id1": 10}},
	}, respbody.Goods)
}

func TestGetGoodsError(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().GetGoods().Return(nil, fmt.Errorf("some error"))
	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(types.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)

	req, err := http.NewRequest(http.MethodGet, s.base+"/api/v1/goods", nil)
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
