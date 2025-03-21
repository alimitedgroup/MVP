package server

import (
	"fmt"
	"net/http"
)

// StartServer avvia un server HTTP di test
func StartServer() {
	/*http.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})*/

	http.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"pong"}`)
	})
	// In server.go
	http.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		// Imposta l'header Content-Type a application/json
		w.Header().Set("Content-Type", "application/json")

		// Restituisci una risposta JSON valida
		fmt.Fprintf(w, `{"token":"token123"}`)
	})
	fmt.Println("Server in ascolto su http://localhost:8080")
	go http.ListenAndServe(":8080", nil)
}
