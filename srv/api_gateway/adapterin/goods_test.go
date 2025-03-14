package adapterin

import (
	"encoding/json"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetGoods(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().GetGoods().Return(
		[]dto.GoodAndAmount{
			{"Apple", "A tasty apple", "id1", 20},
			{"Orange", "A tasty orange", "id2", 10},
		},
		nil,
	)

	resp, err := client.Get(s.base + "/api/v1/goods")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.GetGoodsResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, []dto.GoodAndAmount{
		{"Apple", "A tasty apple", "id1", 20},
		{"Orange", "A tasty orange", "id2", 10},
	}, respbody.Goods)
}

func TestGetGoodsError(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().GetGoods().Return(nil, fmt.Errorf("some error"))

	resp, err := client.Get(s.base + "/api/v1/goods")
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var respbody response.ResponseDTO[any]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "internal_error", respbody.Error)
}
