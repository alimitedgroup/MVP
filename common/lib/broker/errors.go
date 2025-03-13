package broker

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

func RespondToMsg(msg *nats.Msg, resp any) error {
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
