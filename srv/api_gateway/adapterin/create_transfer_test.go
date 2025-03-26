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

func TestCreateTransfer(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.order.EXPECT().CreateTransfer(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)

	payload := bytes.NewReader([]byte(`{"receiver_id": "2", "sender_id": "1", "goods": {"hat-1": 2}}`))
	resp, err := client.Post(s.base+"/api/v1/transfers", "application/json", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.CreateTransferResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.NotEmpty(t, respbody.TransferID)
}
