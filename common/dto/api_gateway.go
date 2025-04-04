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

type AddStockRequest struct {
	WarehouseID string `uri:"warehouse_id"`
	GoodID      string `uri:"good_id"`
	Quantity    int64  `json:"quantity"`
}

type RemoveStockRequest struct {
	WarehouseID string `uri:"warehouse_id"`
	GoodID      string `uri:"good_id"`
	Quantity    int64  `json:"quantity"`
}

type CreateGoodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateGoodRequest struct {
	Id          string `uri:"good_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateGoodResponse struct {
	GoodID string `json:"good_id"`
}

type GetGoodsResponse struct {
	Goods []GoodAndAmount `json:"goods"`
}

type MissingRequiredFieldError struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

type GoodAndAmount struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ID          string           `json:"id"`
	Amount      int64            `json:"amount"`
	Amounts     map[string]int64 `json:"amounts"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
}

type GetOrdersResponse struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	OrderID      string           `json:"order_id"`
	CreationTime int64            `json:"creation_time"`
	UpdateTime   int64            `json:"update_time"`
	Status       string           `json:"status"`
	Name         string           `json:"name"`
	FullName     string           `json:"full_name"`
	Address      string           `json:"address"`
	Goods        map[string]int64 `json:"goods"`
	Reservations []string         `json:"reservations"`
}

type CreateTransferResponse struct {
	TransferID string `json:"transfer_id"`
}

type GetTransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}

type Transfer struct {
	Status       string           `json:"status"`
	CreationTime int64            `json:"creation_time"`
	UpdateTime   int64            `json:"update_time"`
	TransferID   string           `json:"transfer_id"`
	SenderID     string           `json:"sender_id"`
	ReceiverID   string           `json:"receiver_id"`
	Goods        map[string]int64 `json:"goods"`
}

type CreateOrderRequest struct {
	Name     string           `json:"name"`
	FullName string           `json:"full_name"`
	Address  string           `json:"address"`
	Goods    map[string]int64 `json:"goods"`
}

type CreateTransferRequest struct {
	SenderID   string           `json:"sender_id"`
	ReceiverID string           `json:"receiver_id"`
	Goods      map[string]int64 `json:"goods"`
}

type CreateQueryRequest struct {
	GoodID    string `json:"good_id"`
	Operator  string `json:"operator"`
	Threshold int    `json:"threshold"`
}

type CreateQueryResponse struct {
	QueryID string `json:"query_id"`
}

type GetQueriesResponse struct {
	Queries []Query `json:"queries"`
}

type Query struct {
	QueryID   string `json:"query_id"`
	GoodID    string `json:"good_id"`
	Operator  string `json:"operator"`
	Threshold int    `json:"threshold"`
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
