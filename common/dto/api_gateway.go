package dto

import (
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto/response"
)

type AuthLoginRequest struct {
	Username string `json:"username"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}

type IsLoggedResponse struct {
	Role string `json:"role"`
}

type GetWarehousesResponse struct {
	Ids []string `json:"warehouse_ids"`
}

type GetGoodsResponse struct {
	Goods []GoodAndAmount `json:"goods"`
}

type MissingRequiredFieldError struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

type GoodAndAmount struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Amount      int64  `json:"amount"`
}

func FieldIsRequired(fieldName string) response.ResponseDTO[MissingRequiredFieldError] {
	return response.ResponseDTO[MissingRequiredFieldError]{
		Error: "missing_field",
		Message: MissingRequiredFieldError{
			Field:       fieldName,
			Description: fmt.Sprintf("The `%s` field is required", fieldName),
		},
	}
}

func InternalError() response.ResponseDTO[string] {
	return response.ResponseDTO[string]{
		Error:   "internal_error",
		Message: "No further details are available",
	}
}

func AuthFailed() response.ResponseDTO[string] {
	return response.ResponseDTO[string]{
		Error:   "authentication_failed",
		Message: "The credentials you provided are invalid",
	}
}

func MissingToken() response.ResponseDTO[string] {
	return response.ResponseDTO[string]{
		Error:   "missing_token",
		Message: "You didn't provide a token in your request. Refer to the manual for more information",
	}
}

func InvalidToken() response.ResponseDTO[string] {
	return response.ResponseDTO[string]{
		Error:   "invalid_token",
		Message: "The token you provided is invalid. Refer to the manual for more information",
	}
}

func ExpiredToken() response.ResponseDTO[string] {
	return response.ResponseDTO[string]{
		Error:   "expired_token",
		Message: "The token you provided is expired. You should login again.",
	}
}
