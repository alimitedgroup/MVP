/*

PER PROVE CON SERVER- DA ELIMINARE

package server

import (
	"fmt"
	"net/http"
)

// StartServer avvia un server HTTP di test
func StartServer() {
	http.HandleFunc("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	http.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "token123")
	})
	fmt.Println("Server in ascolto su http://localhost:8080")
	go http.ListenAndServe(":8080", nil)
}
*/