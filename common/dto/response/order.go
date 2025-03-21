package response

type OrderCreateResponseDTO ResponseDTO[OrderCreateInfo]

type OrderCreateInfo struct {
	OrderID string `json:"order_id"`
}

type GetOrderResponseDTO ResponseDTO[OrderInfo]
type GetAllOrderResponseDTO ResponseDTO[[]OrderInfo]

type OrderInfo struct {
	OrderID      string          `json:"order_id"`
	Status       string          `json:"status"`
	Name         string          `json:"name"`
	FullName     string          `json:"full_name"`
	Address      string          `json:"address"`
	Goods        []OrderInfoGood `json:"goods"`
	Reservations []string        `json:"reservations"`
}

type OrderInfoGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type TransferCreateResponseDTO ResponseDTO[TransferCreateInfo]
type GetTransferResponseDTO ResponseDTO[TransferInfo]
type GetAllTransferResponseDTO ResponseDTO[[]TransferInfo]

type TransferCreateInfo struct {
	TransferID string `json:"transfer_id"`
}

type TransferInfo struct {
	Status     string             `json:"status"`
	TransferID string             `json:"transfer_id"`
	SenderID   string             `json:"sender_id"`
	ReceiverID string             `json:"receiver_id"`
	Goods      []TransferInfoGood `json:"goods"`
}

type TransferInfoGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}
