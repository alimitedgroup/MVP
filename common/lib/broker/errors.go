package broker

import (
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/nats-io/nats.go"
)

func RespondToMsg[T any](msg *nats.Msg, resp response.ResponseDTO[T]) error {
	payload, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	err = msg.Respond(payload)
	if err != nil {
		return err
	}

	return nil
}
