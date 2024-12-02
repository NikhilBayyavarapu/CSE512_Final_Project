package main

import (
	"cse512/db"
	"cse512/handlers"
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	port := flag.Int("port", 0, "Port to run the server on")
	help := flag.Bool("help", false, "Use port flag to specify port to run the server on")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	if *port == 0 {
		fmt.Println("Please specify a port to run the server on using -port flag")
		return
	}

	_ = db.GetClient()
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.HandleLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/transactions", handlers.HandleTransaction).Methods("GET", "OPTIONS")
	router.HandleFunc("/transaction", handlers.PerformTransaction).Methods("POST", "OPTIONS")
	router.HandleFunc("/monthdata", handlers.GetMonthData).Methods("GET", "OPTIONS")

	fmt.Printf("Starting server on port %d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
