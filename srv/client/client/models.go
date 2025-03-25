package client

// Response structures for different API endpoints
type PingResponse struct {
	Message string `json:"message"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}

type IsLoggedResponse struct {
	Role string `json:"role"`
}

type GetWarehousesResponse struct {
	Ids []string `json:"ids"`
}

type GoodAndAmount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int64  `json:"quantity"`
}

type GetGoodsResponse struct {
	Goods []GoodAndAmount `json:"goods"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}
