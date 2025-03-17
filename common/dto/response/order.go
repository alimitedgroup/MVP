package response

type OrderCreateResponseDTO ResponseDTO[OrderCreateInfo]

type OrderCreateInfo struct {
	OrderID string `json:"order_id"`
}

type ErrorResponseDTO ResponseDTO[any]

type GetOrderResponseDTO ResponseDTO[OrderInfo]
type GetAllOrderResponseDTO ResponseDTO[[]OrderInfo]

type OrderInfo struct {
	OrderID      string          `json:"order_id"`
	Status       string          `json:"status"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	Address      string          `json:"address"`
	Goods        []OrderInfoGood `json:"goods"`
	Reservations []string        `json:"reservations"`
}

type OrderInfoGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type TransferInfo struct {
	TransferID string             `json:"transfer_id"`
	SenderID   string             `json:"sender_id"`
	ReceiverID string             `json:"receiver_id"`
	Goods      []TransferInfoGood `json:"goods"`
}

type TransferInfoGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}
