/*package main

import (
	"fmt"
	"time"

	"go.uber.org/fx"
)

// Example application using the client
func main() {
	app := fx.New(
		// Include our API client module
		Module,

		// Override default config if needed
		fx.Decorate(func() ClientConfig {
			return ClientConfig{
				BaseURL: "https://api.example.com",
				Timeout: time.Second * 60,
			}
		}),

		// Use the client in the application
		fx.Invoke(func(client *Client) {
			// Use the client here...
			ping, err := client.Ping()
			if err != nil {
				fmt.Printf("Error pinging API: %v\n", err)
				return
			}
			fmt.Printf("API response: %s\n", ping.Message)

			// Login and use authenticated endpoints
			auth, err := client.Login("testuser")
			if err != nil {
				fmt.Printf("Error logging in: %v\n", err)
				return
			}
			fmt.Printf("Logged in with token: %s\n", auth.Token)

			// Get user role
			logged, err := client.IsLogged()
			if err != nil {
				fmt.Printf("Error checking login: %v\n", err)
				return
			}
			fmt.Printf("User role: %s\n", logged.Role)
		}),
	)

	app.Run()
}
*/

package main

import (
	"fmt"
	"github.com/alimitedgroup/MVP/srv/client/client"
	"go.uber.org/fx"
	"time"
)

func main() {
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
