package main

import (
	"fmt"
	"time"

	"github.com/alimitedgroup/MVP/srv/client/client"
	"github.com/alimitedgroup/MVP/srv/client/server"
	"go.uber.org/fx"
)

func main() {

	// PER PROVE CON SERVER
	server.StartServer()

	// Breve pausa per dare tempo al server di avviarsi
	time.Sleep(time.Second)

	app := fx.New(
		fx.Provide(
			// Fornisci ClientConfig
			func() client.ClientConfig {
				return client.ClientConfig{
					BaseURL: "http://localhost:8080",
					Timeout: time.Second * 30,
				}
			},
			// Fornisci il client
			client.NewClient,
		),
		fx.Invoke(func(apiClient *client.Client) {
			// Test Ping
			pingMessage, err := apiClient.Ping()
			if err != nil {
				fmt.Printf("Error pinging API: %v\n", err)
				return
			}
			fmt.Printf("API response: %s\n", pingMessage)

			// Test Login
			token, err := apiClient.Login("testuser")
			if err != nil {
				fmt.Printf("Error logging in: %v\n", err)
				return
			}
			fmt.Printf("Logged in with token: %s\n", token)
		}),
	)

	app.Run()
}
