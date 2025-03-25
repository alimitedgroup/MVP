package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimitedgroup/MVP/srv/client/client"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	baseURL string
	timeout int
)

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Use:   "apicli",
		Short: "API Client CLI Tool",
		Long:  "Command line interface for interacting with the warehouse API",
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "http://localhost:8080", "API base URL")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "Request timeout in seconds")

	// Ping command
	pingCmd := &cobra.Command{
		Use:   "ping",
		Short: "Check if API server is running",
		Run: func(cmd *cobra.Command, args []string) {
			app := fx.New(
				fx.NopLogger,
				fx.Provide(
					func() client.ClientConfig {
						return client.ClientConfig{
							BaseURL: baseURL,
							Timeout: time.Duration(timeout) * time.Second,
						}
					},
					client.NewClient,
				),
				fx.Invoke(func(apiClient *client.Client) {
					pingResponse, err := apiClient.Ping()
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						return
					}
					fmt.Printf("Server response: %s\n", pingResponse.Message)
				}),
			)
			app.Start(cmd.Context())
			defer app.Stop(cmd.Context())
		},
	}

	// Login command
	loginCmd := &cobra.Command{
		Use:   "login [username]",
		Short: "Authenticate and get a token",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := ""
			if len(args) > 0 {
				username = args[0]
			} else {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter username: ")
				input, _ := reader.ReadString('\n')
				username = strings.TrimSpace(input)
			}

			app := fx.New(
				fx.NopLogger,
				fx.Provide(
					func() client.ClientConfig {
						return client.ClientConfig{
							BaseURL: baseURL,
							Timeout: time.Duration(timeout) * time.Second,
						}
					},
					client.NewClient,
				),
				fx.Invoke(func(apiClient *client.Client) {
					loginResponse, err := apiClient.Login(username)
					if err != nil {
						fmt.Printf("Login failed: %v\n", err)
						return
					}
					fmt.Printf("Login successful\n")
					fmt.Printf("Token: %s\n", loginResponse.Token)

					// Save token to file for other commands to use
					saveToken(loginResponse.Token)
				}),
			)
			app.Start(cmd.Context())
			defer app.Stop(cmd.Context())
		},
	}

	// IsLogged command
	isLoggedCmd := &cobra.Command{
		Use:   "status",
		Short: "Check login status and user role",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := loadToken()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			app := fx.New(
				fx.NopLogger,
				fx.Provide(
					func() client.ClientConfig {
						return client.ClientConfig{
							BaseURL: baseURL,
							Timeout: time.Duration(timeout) * time.Second,
						}
					},
					client.NewClient,
				),
				fx.Invoke(func(apiClient *client.Client) {
					apiClient.Token = token
					statusResponse, err := apiClient.IsLogged()
					if err != nil {
						fmt.Printf("Error checking status: %v\n", err)
						return
					}
					fmt.Printf("Logged in as: %s\n", statusResponse.Role)
				}),
			)
			app.Start(cmd.Context())
			defer app.Stop(cmd.Context())
		},
	}

	// GetWarehouses command
	getWarehousesCmd := &cobra.Command{
		Use:   "warehouses",
		Short: "List all warehouses",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := loadToken()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			app := fx.New(
				fx.NopLogger,
				fx.Provide(
					func() client.ClientConfig {
						return client.ClientConfig{
							BaseURL: baseURL,
							Timeout: time.Duration(timeout) * time.Second,
						}
					},
					client.NewClient,
				),
				fx.Invoke(func(apiClient *client.Client) {
					apiClient.Token = token
					warehousesResponse, err := apiClient.GetWarehouses()
					if err != nil {
						fmt.Printf("Error fetching warehouses: %v\n", err)
						return
					}

					fmt.Println("Warehouses:")
					for i, id := range warehousesResponse.Ids {
						fmt.Printf("%d. %s\n", i+1, id)
					}
				}),
			)
			app.Start(cmd.Context())
			defer app.Stop(cmd.Context())
		},
	}

	// GetGoods command
	getGoodsCmd := &cobra.Command{
		Use:   "goods",
		Short: "List all goods with quantities",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := loadToken()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			app := fx.New(
				fx.NopLogger,
				fx.Provide(
					func() client.ClientConfig {
						return client.ClientConfig{
							BaseURL: baseURL,
							Timeout: time.Duration(timeout) * time.Second,
						}
					},
					client.NewClient,
				),
				fx.Invoke(func(apiClient *client.Client) {
					apiClient.Token = token
					goodsResponse, err := apiClient.GetGoods()
					if err != nil {
						fmt.Printf("Error fetching goods: %v\n", err)
						return
					}

					fmt.Println("Goods inventory:")
					fmt.Printf("%-10s %-30s %-10s %s\n", "ID", "Name", "Amount", "Description")
					fmt.Println(strings.Repeat("-", 80))

					for _, good := range goodsResponse.Goods {
						fmt.Printf("%-10s %-30s %-10d %s\n",
							good.ID,
							truncateString(good.Name, 30),
							good.Quantity,
							truncateString(good.Description, 30))
					}
				}),
			)
			app.Start(cmd.Context())
			defer app.Stop(cmd.Context())
		},
	}

	// Interactive mode command
	interactiveCmd := &cobra.Command{
		Use:   "interactive",
		Short: "Start interactive mode",
		Run: func(cmd *cobra.Command, args []string) {
			runInteractiveMode()
		},
	}

	// Add all commands to root
	rootCmd.AddCommand(pingCmd, loginCmd, isLoggedCmd, getWarehousesCmd, getGoodsCmd, interactiveCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper functions

// saveToken saves the token to a file in the user's home directory
func saveToken(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	tokenFile := homeDir + "/.apicli_token"
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

// loadToken loads the token from a file in the user's home directory
func loadToken() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	tokenFile := homeDir + "/.apicli_token"
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("not logged in, please run 'login' first")
	}

	return string(data), nil
}

// truncateString truncates a string to the specified length and adds ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// runInteractiveMode starts an interactive CLI session
func runInteractiveMode() {
	fmt.Println("Interactive Mode - Type 'help' for available commands, 'exit' to quit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		command := args[0]
		cmdArgs := args[1:]

		switch command {
		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  ping               - Check if API server is running")
			fmt.Println("  login [username]   - Authenticate and get a token")
			fmt.Println("  status             - Check login status and user role")
			fmt.Println("  warehouses         - List all warehouses")
			fmt.Println("  goods              - List all goods with quantities")
			fmt.Println("  exit               - Exit interactive mode")

		case "ping":
			os.Args = []string{"apicli", "ping"}
			pingCmd := &cobra.Command{Use: "ping"}
			pingCmd.Run(pingCmd, cmdArgs)

		case "login":
			os.Args = append([]string{"apicli", "login"}, cmdArgs...)
			loginCmd := &cobra.Command{Use: "login"}
			loginCmd.Run(loginCmd, cmdArgs)

		case "status":
			os.Args = []string{"apicli", "status"}
			statusCmd := &cobra.Command{Use: "status"}
			statusCmd.Run(statusCmd, cmdArgs)

		case "warehouses":
			os.Args = []string{"apicli", "warehouses"}
			warehousesCmd := &cobra.Command{Use: "warehouses"}
			warehousesCmd.Run(warehousesCmd, cmdArgs)

		case "goods":
			os.Args = []string{"apicli", "goods"}
			goodsCmd := &cobra.Command{Use: "goods"}
			goodsCmd.Run(goodsCmd, cmdArgs)

		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}
	}
}
