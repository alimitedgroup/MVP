package request

type AddStockRequestDTO struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type RemoveStockRequestDTO struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}
