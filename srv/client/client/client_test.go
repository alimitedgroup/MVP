package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/ping", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"pong"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		BaseURL: mockServer.URL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	resp, err := client.Ping()
	assert.NoError(t, err)
	assert.Equal(t, "pong", resp.Message)
}

func TestLogin(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/login", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":"fake-token"}`))
	}))
	defer mockServer.Close()

	client := &Client{
		BaseURL: mockServer.URL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	resp, err := client.Login("user")
	assert.NoError(t, err)
	assert.Equal(t, "fake-token", resp.Token)
}

func TestGetWarehouses(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/warehouses", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ids":["Main Warehouse"]}`))
	}))
	defer mockServer.Close()

	client := &Client{
		BaseURL: mockServer.URL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	warehouses, err := client.GetWarehouses()
	assert.NoError(t, err)
	assert.Len(t, warehouses.Ids, 1)
	assert.Equal(t, "Main Warehouse", warehouses.Ids[0])
}

func TestGetGoods(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/goods", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"goods":[{"id":"1","name":"Widget","description":"A test widget","amount":100}]}`))
	}))
	defer mockServer.Close()

	client := &Client{
		BaseURL: mockServer.URL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	goods, err := client.GetGoods()
	assert.NoError(t, err)
	assert.Len(t, goods.Goods, 1)
	assert.Equal(t, "Widget", goods.Goods[0].Name)
	assert.Equal(t, int64(100), goods.Goods[0].Amount)
}
