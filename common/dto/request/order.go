package request

type CreateOrderRequestDTO struct {
	Name     string            `json:"name"`
	FullName string            `json:"full_name"`
	Address  string            `json:"address"`
	Goods    []CreateOrderGood `json:"goods"`
}

type CreateOrderGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type CreateTransferRequestDTO struct {
	SenderID   string         `json:"sender_id"`
	ReceiverID string         `json:"receiver_id"`
	Goods      []TransferGood `json:"goods"`
}

type TransferGood struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type GetOrderRequestDTO struct {
	OrderID string `json:"order_id"`
}

type GetTransferRequestDTO struct {
	TransferID string `json:"transfer_id"`
}
