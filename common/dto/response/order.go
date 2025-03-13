package response

type OrderCreateResponseDTO ResponseDTO[OrderCreateInfo]

type OrderCreateInfo struct {
	OrderID string `json:"order_id"`
}

type ErrorResponseDTO ResponseDTO[any]
