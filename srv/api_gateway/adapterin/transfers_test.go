package adapterin

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/stretchr/testify/require"
)

func TestGetTransfers(t *testing.T) {
	s := start(t)
	client := &http.Client{}

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

	resp, err := client.Get(s.base + "/api/v1/transfers")
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
