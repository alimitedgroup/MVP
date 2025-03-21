package request

type AddStockRequestDTO struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type RemoveStockRequestDTO struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}

type ReserveStockRequestDTO struct {
	Goods []ReserveStockItem `json:"goods"`
}

type ReserveStockItem struct {
	GoodID   string `json:"good_id"`
	Quantity int64  `json:"quantity"`
}
