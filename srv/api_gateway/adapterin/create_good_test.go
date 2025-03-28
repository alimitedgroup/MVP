package adapterin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateGood(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(portin.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)
	s.warehouses.EXPECT().CreateGood(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)

	payload := bytes.NewReader([]byte(`{"name":"hat","description":"blue hat"}`))
	req, err := http.NewRequest(http.MethodPost, s.base+"/api/v1/goods", payload)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.CreateGoodResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "1", respbody.GoodID)
}
