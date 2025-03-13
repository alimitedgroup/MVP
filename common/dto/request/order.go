package request

type CreateOrderRequestDTO struct {
	Name    string            `json:"name"`
	Email   string            `json:"email"`
	Address string            `json:"address"`
	Goods   []CreateOrderGood `json:"goods"`
}

type CreateOrderGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type TransferRequestDTO struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
}
