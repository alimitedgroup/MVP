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

// ---- Strutture per le risposte dell'API ----

// AuthResponse rappresenta la risposta dell'endpoint di autenticazione
type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// TokenValidateResponse rappresenta la risposta della validazione token
type TokenValidateResponse struct {
	Valid bool   `json:"valid"`
	Role  string `json:"role,omitempty"`
}

// PingResponse rappresenta la risposta dell'endpoint di ping
type PingResponse struct {
	Message string `json:"message"`
}

// IsLoggedResponse rappresenta la risposta dell'endpoint di verifica login
type IsLoggedResponse struct {
	Role string `json:"role"`
}

// GetWarehousesResponse rappresenta la risposta dell'endpoint dei magazzini
type GetWarehousesResponse struct {
	Ids []string `json:"ids"`
}

// GoodAndAmount rappresenta un prodotto con la sua quantità
type GoodAndAmount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
}

// GetGoodsResponse rappresenta la risposta dell'endpoint dei prodotti
type GetGoodsResponse struct {
	Goods []GoodAndAmount `json:"goods"`
}

// ErrorResponse rappresenta una risposta di errore dall'API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ---- Configurazione del client ----

// APIClientConfig contiene la configurazione per l'APIClient
type APIClientConfig struct {
	BaseURL               string
	AuthServiceURL        string
	Timeout               time.Duration
	TokenRefreshThreshold time.Duration // Tempo prima della scadenza per refreshare il token
}

// APIClientParams definisce i parametri di dependency injection per costruire un APIClient
type APIClientParams struct {
	fx.In

	Config APIClientConfig
}

// APIClient rappresenta il client per l'API e il servizio di autenticazione
type APIClient struct {
	BaseURL               string
	AuthServiceURL        string
	HTTPClient            *http.Client
	Token                 string
	TokenExpiry           time.Time
	TokenRefreshThreshold time.Duration
}

// ProvideAPIClientConfig crea una configurazione predefinita per APIClient
func ProvideAPIClientConfig() APIClientConfig {
	return APIClientConfig{
		BaseURL:               "http://api-service:8080",
		AuthServiceURL:        "http://authenticator-service:8080",
		Timeout:               time.Second * 30,
		TokenRefreshThreshold: time.Hour * 24, // Refresha se manca meno di 1 giorno alla scadenza
	}
}

// NewAPIClient crea un nuovo APIClient con dependency injection
func NewAPIClient(p APIClientParams) *APIClient {
	return &APIClient{
		BaseURL:        p.Config.BaseURL,
		AuthServiceURL: p.Config.AuthServiceURL,
		HTTPClient: &http.Client{
			Timeout: p.Config.Timeout,
		},
		TokenRefreshThreshold: p.Config.TokenRefreshThreshold,
	}
}

// APIModule fornisce i componenti fx per APIClient
var APIModule = fx.Options(
	fx.Provide(ProvideAPIClientConfig),
	fx.Provide(NewAPIClient),
)

// ---- Metodi di autenticazione ----

// Login autentica un utente e ottiene un token
func (c *APIClient) Login(username, password string) error {
	// Prepara i dati per la richiesta di login
	data := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("errore nella codifica JSON dei dati di login: %w", err)
	}

	req, err := http.NewRequest("POST", c.AuthServiceURL+"/api/v1/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("errore nella creazione della richiesta di login: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var authResp AuthResponse
	if err := c.sendRequest(req, &authResp); err != nil {
		return fmt.Errorf("errore nell'autenticazione: %w", err)
	}

	// Salva il token e la sua scadenza
	c.Token = authResp.Token
	c.TokenExpiry = authResp.ExpiresAt

	return nil
}

// LoginWithUsernameForm autentica un utente utilizzando form data invece di JSON
// Questo è un metodo alternativo per supportare l'API esistente
func (c *APIClient) LoginWithUsernameForm(username string) error {
	data := bytes.NewBufferString(fmt.Sprintf("username=%s", username))
	req, err := http.NewRequest("POST", c.BaseURL+"/api/v1/login", data)
	if err != nil {
		return fmt.Errorf("errore nella creazione della richiesta: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var authResp struct {
		Token string `json:"token"`
	}
	if err := c.sendRequest(req, &authResp); err != nil {
		return fmt.Errorf("errore nell'autenticazione: %w", err)
	}

	// Salva il token (qui non abbiamo informazioni sulla scadenza)
	c.Token = authResp.Token
	// Imposta una scadenza predefinita di 1 settimana come da specifica
	c.TokenExpiry = time.Now().Add(7 * 24 * time.Hour)

	return nil
}

// ValidateToken verifica se il token corrente è valido
func (c *APIClient) ValidateToken() (*TokenValidateResponse, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("nessun token disponibile, effettuare il login")
	}

	req, err := http.NewRequest("GET", c.AuthServiceURL+"/api/v1/auth/validate", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	res := TokenValidateResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, fmt.Errorf("errore nella validazione del token: %w", err)
	}

	return &res, nil
}

// RefreshTokenIfNeeded rinnova il token se è vicino alla scadenza
func (c *APIClient) RefreshTokenIfNeeded() error {
	// Se non c'è un token, non c'è niente da refreshare
	if c.Token == "" {
		return fmt.Errorf("nessun token disponibile, effettuare il login")
	}

	// Controlla se il token è vicino alla scadenza
	if time.Until(c.TokenExpiry) > c.TokenRefreshThreshold {
		// Il token è ancora valido per abbastanza tempo
		return nil
	}

	// Richiedi un nuovo token usando il token corrente
	req, err := http.NewRequest("POST", c.AuthServiceURL+"/api/v1/auth/refresh", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	var authResp AuthResponse
	if err := c.sendRequest(req, &authResp); err != nil {
		return fmt.Errorf("errore nel refresh del token: %w", err)
	}

	// Aggiorna il token e la sua scadenza
	c.Token = authResp.Token
	c.TokenExpiry = authResp.ExpiresAt

	return nil
}

// ---- Metodi dell'API generale ----

// Ping verifica se il server API è online
func (c *APIClient) Ping() (*PingResponse, error) {
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

// IsLogged verifica se l'utente è autenticato e ottiene il ruolo
func (c *APIClient) IsLogged() (*IsLoggedResponse, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("non autenticato, effettuare prima il login")
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

// GetWarehouses ottiene la lista dei magazzini
func (c *APIClient) GetWarehouses() (*GetWarehousesResponse, error) {
	// Aggiorna il token se necessario
	if err := c.RefreshTokenIfNeeded(); err != nil {
		return nil, fmt.Errorf("errore nel refresh del token: %w", err)
	}

	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/warehouses", nil)
	if err != nil {
		return nil, err
	}

	// Aggiungi il token di autorizzazione se disponibile
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	res := GetWarehousesResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetGoods ottiene la lista dei prodotti con le relative quantità
func (c *APIClient) GetGoods() (*GetGoodsResponse, error) {
	// Aggiorna il token se necessario
	if err := c.RefreshTokenIfNeeded(); err != nil {
		return nil, fmt.Errorf("errore nel refresh del token: %w", err)
	}

	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/goods", nil)
	if err != nil {
		return nil, err
	}

	// Aggiungi il token di autorizzazione se disponibile
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	res := GetGoodsResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ---- Metodi di supporto ----

// sendRequest invia la richiesta e deserializza la risposta
func (c *APIClient) sendRequest(req *http.Request, v interface{}) error {
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
		// Prova a deserializzare una risposta di errore
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			return fmt.Errorf("errore API: %s (codice: %s)", errorResponse.Message, errorResponse.Code)
		}
		return fmt.Errorf("richiesta API fallita con status code %d: %s", resp.StatusCode, body)
	}

	return json.Unmarshal(body, v)
}
