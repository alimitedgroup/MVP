package adapterin

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/stretchr/testify/require"
)

func TestGetTransfers(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.auth.EXPECT().ValidateToken("some.secure.jwt").Return(portin.UserData{
		Username: "test",
		Role:     types.RoleGlobalAdmin,
	}, nil)

	s.order.EXPECT().GetTransfers().Return(
		[]dto.Transfer{
			{
				TransferID: "1",
				Status:     "Created",
				SenderID:   "1",
				ReceiverID: "2",
				Goods:      map[string]int64{"1": 1},
			},
		},
		nil,
	)

	req, err := http.NewRequest(http.MethodGet, s.base+"/api/v1/transfers", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer some.secure.jwt")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.GetTransfersResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, []dto.Transfer{
		{
			TransferID: "1",
			Status:     "Created",
			SenderID:   "1",
			ReceiverID: "2",
			Goods:      map[string]int64{"1": 1},
		},
	}, respbody.Transfers)
}
