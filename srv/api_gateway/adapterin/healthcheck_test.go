package adapterin

import (
	"encoding/json"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	resp, err := client.Get(s.base + "/api/v1/ping")
	require.NoError(t, err)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)

	require.Equal(t, "pong", respbody.Message)
	require.Zero(t, respbody.Error)
}
