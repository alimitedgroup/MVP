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

func TestCreateGood(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().CreateGood(gomock.Any(), gomock.Any(), gomock.Any()).Return("1", nil)

	payload := bytes.NewReader([]byte(`{"name":"hat","description":"blue hat"}`))
	resp, err := client.Post(s.base+"/api/v1/goods", "application/json", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody dto.CreateGoodResponse
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "1", respbody.GoodID)
}
