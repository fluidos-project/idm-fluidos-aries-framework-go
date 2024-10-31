package main

import (
	"fabricrest-go/api-rest/web"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Model Treatment endpoints
	router.HandleFunc("/model/average", web.CalculateAverageModelUpdate).Methods("POST")
	router.HandleFunc("/model/{id}", web.ReadAverageModelUpdate).Methods("GET")
	router.HandleFunc("/models", web.GetAllModelUpdates).Methods("GET")
	router.HandleFunc("/models/query", web.QueryModelUpdatesByDateRange).Methods("GET")

	fmt.Printf("Starting server on http://localhost:3005/\n")
	if err := http.ListenAndServe(":3005", router); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
