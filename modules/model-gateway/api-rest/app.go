package main

import (
	"fabricrest-go/api-rest/web"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Existing aggregation endpoints
	router.HandleFunc("/model/aggregated", web.AggregateModel).Methods("POST")
	router.HandleFunc("/model/aggregated/{id}", web.ReadAggregatedModel).Methods("GET")
	router.HandleFunc("/models/aggregated", web.GetAllAggregatedModels).Methods("GET")
	router.HandleFunc("/models/aggregated/query", web.QueryAggregatedModelsByDateRange).Methods("GET")

	// New base model update endpoints
	router.HandleFunc("/basemodel/update", web.CreateBaseModelUpdate).Methods("POST")
	router.HandleFunc("/basemodel/{id}", web.ReadBaseModelUpdate).Methods("GET")
	router.HandleFunc("/basemodels", web.GetAllBaseModelUpdates).Methods("GET")
	router.HandleFunc("/basemodels/query", web.QueryBaseModelUpdatesByDateRange).Methods("GET")
	fmt.Printf("Starting server on http://localhost:3005/\n")
	if err := http.ListenAndServe(":3005", router); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
