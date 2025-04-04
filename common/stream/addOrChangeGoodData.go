package stream

import "github.com/nats-io/nats.go/jetstream"

var AddOrChangeGoodDataStream = jetstream.StreamConfig{
	Name:     "good_data_update",
	Subjects: []string{"good.update"},
	Storage:  jetstream.FileStorage,
}

type GoodUpdateData struct {
	GoodID             string `json:"id"` //must be a new goodID if a new good was added, otherwise must be a valid goodID
	GoodNewName        string `json:"new_name"`
	GoodNewDescription string `json:"new_description"`
}
