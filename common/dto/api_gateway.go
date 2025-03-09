package dto

type AuthLoginRequest struct {
	Username string `json:"username"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}
