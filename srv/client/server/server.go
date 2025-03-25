package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StartServer avvia un server HTTP di test
func StartServer() {
	// Existing ping endpoint
	http.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	// Login endpoint
	http.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": "token123"})
	})

	// Add warehouses endpoint
	http.HandleFunc("/api/v1/warehouses", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{
			"ids": {"Warehouse1", "Warehouse2", "Warehouse3"},
		})
	})

	// Add goods endpoint
	http.HandleFunc("/api/v1/goods", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]map[string]interface{}{
			"goods": {
				{
					"id":          "1",
					"name":        "Widget",
					"description": "A test widget",
					"amount":      100,
				},
				{
					"id":          "2",
					"name":        "Gadget",
					"description": "Another test item",
					"amount":      50,
				},
			},
		})
	})

	fmt.Println("Server in ascolto su http://localhost:8080")
	go http.ListenAndServe(":8080", nil)
}
