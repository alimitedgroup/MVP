package stream

type AddQueryRule struct {
	GoodID    string `json:"good_id"`
	Operator  string `json:"operator"`
	Threshold int    `json:"threshold"`
}
