package adapterin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestAddStock(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().AddStock(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	payload := bytes.NewReader([]byte(`{"quantity": 10}`))
	req, err := http.NewRequest(http.MethodPost, s.base+"/api/v1/goods/1/warehouse/1/add", payload)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respbody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Empty(t, respbody)
}

func TestAddStockError(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	s.warehouses.EXPECT().AddStock(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))

	payload := bytes.NewReader([]byte(`{"quantity": 10}`))

	req, err := http.NewRequest(http.MethodPost, s.base+"/api/v1/goods/1/warehouse/1/add", payload)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "internal_error", respbody.Error)
}

func TestAddStockMalformedInput(t *testing.T) {
	s := start(t)
	client := &http.Client{}

	payload := bytes.NewReader([]byte(`"quantity": 10`))

	req, err := http.NewRequest(http.MethodPost, s.base+"/api/v1/goods/1/warehouse/1/add", payload)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var respbody response.ResponseDTO[string]
	err = json.NewDecoder(resp.Body).Decode(&respbody)
	require.NoError(t, err)
	require.Equal(t, "internal_error", respbody.Error)
}
