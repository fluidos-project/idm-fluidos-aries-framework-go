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
	router.HandleFunc("/models/aggregate", web.AggregateModel).Methods("POST")
	router.HandleFunc("/models/aggregate/{id}", web.ReadAggregatedModel).Methods("GET")
	router.HandleFunc("/models/aggregate", web.GetAllAggregatedModels).Methods("GET")
	router.HandleFunc("/models/aggregate/query", web.QueryAggregatedModelsByDateRange).Methods("GET")

	fmt.Printf("Starting server on http://localhost:3005/\n")
	if err := http.ListenAndServe(":3005", router); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
