package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/fx"
)

// ClientConfig holds configuration for the API client
type ClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

// ClientParams defines dependency injection parameters for constructing a Client
type ClientParams struct {
	fx.In

	Config ClientConfig
}

// Client represents the API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new API client with dependency injection
func NewClient(p ClientParams) *Client {
	return &Client{
		BaseURL: p.Config.BaseURL,
		HTTPClient: &http.Client{
			Timeout: p.Config.Timeout,
		},
	}
}

// ProvideClientConfig creates a ClientConfig with default values
func ProvideClientConfig() ClientConfig {
	return ClientConfig{
		BaseURL: "http://localhost:8080",
		Timeout: time.Second * 30,
	}
}

// Module provides fx components for API client
var Module = fx.Options(
	fx.Provide(ProvideClientConfig),
	fx.Provide(NewClient),
)

// Ping checks if the API server is up and running
func (c *Client) Ping() (*PingResponse, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/ping", nil)
	if err != nil {
		return nil, err
	}

	res := PingResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Login authenticates a user and retrieves a token
func (c *Client) Login(username string) (*AuthLoginResponse, error) {
	data := bytes.NewBufferString(fmt.Sprintf("username=%s", username))
	req, err := http.NewRequest("POST", c.BaseURL+"/api/v1/login", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := AuthLoginResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	// Store the token for future authenticated requests
	c.Token = res.Token

	return &res, nil
}

// IsLogged checks if the current token is valid and returns user role
func (c *Client) IsLogged() (*IsLoggedResponse, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("not authenticated, please login first")
	}

	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/is_logged", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	res := IsLoggedResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetWarehouses lists all warehouses
func (c *Client) GetWarehouses() (*GetWarehousesResponse, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/warehouses", nil)
	if err != nil {
		return nil, err
	}

	res := GetWarehousesResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetGoods lists all goods with their quantities
func (c *Client) GetGoods() (*GetGoodsResponse, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/goods", nil)
	if err != nil {
		return nil, err
	}

	res := GetGoodsResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// sendRequest does the actual request and unmarshals the response
func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		// Try to unmarshal error response
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			return fmt.Errorf("API error: %s (code: %s)", errorResponse.Message, errorResponse.Code)
		}
		return fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, body)
	}

	return json.Unmarshal(body, v)
}
