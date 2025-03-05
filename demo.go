package main

import (
	"encoding/json"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("auth.login", func(msg *nats.Msg) {
		var req dto.AuthLoginRequest
		json.Unmarshal(msg.Data, &req)
		if req.Username == "admin" {
			jsonBoh, err := json.Marshal(dto.AuthLoginResponse{"abc123abc123"})
			if err != nil {
				panic(err)
			}
			err = msg.Respond(jsonBoh)
			if err != nil {
				panic(err)
			}
		} else {
			jsonBoh, _ := json.Marshal(response.ResponseDTO[string]{Error: "invalid_credentials", Message: "Invalid credentials"})
			err = msg.Respond(jsonBoh)
			if err != nil {
				panic(err)
			}
		}

	})
	if err != nil {
		panic(err)
	}

	select {}
}
