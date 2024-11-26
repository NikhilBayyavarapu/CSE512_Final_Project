package main

import (
	"cse512/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.HandleLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/transactions", handlers.HandleTransaction).Methods("GET", "OPTIONS")

	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
