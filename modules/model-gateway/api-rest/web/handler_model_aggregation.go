package web

import (
	"encoding/json"
	application_gateway "fabricrest-go/application-gateway"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type AggregateModelRequest struct {
	Data             [][]float64 `json:"data"`
	BaseModel        string      `json:"baseModel"`
	BaseModelVersion string      `json:"baseModelVersion"`
	Date             string      `json:"date"`
	NodeDID          string      `json:"nodeDID"`
	SignedProof      string      `json:"signedProof"`
}

type CalculationResponse struct {
	Status    string   `json:"status"`
	Message   string   `json:"message"`
	TxID      string   `json:"transactionId"`
	ID        string   `json:"id"`
	ModelsRef []string `json:"modelsRef"`
}

// AggregateModel handles the calculation and storage of aggregate model updates
func AggregateModel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - AggregateModel")

	var request AggregateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if request.Data == nil {
		http.Error(w, "Data matrix cannot be nil", http.StatusBadRequest)
		return
	}
	if request.BaseModel == "" {
		http.Error(w, "BaseModel is required", http.StatusBadRequest)
		return
	}
	if request.BaseModelVersion == "" {
		http.Error(w, "BaseModelVersion is required", http.StatusBadRequest)
		return
	}

	result, err := application_gateway.AggregateModel(
		request.Data,
		request.BaseModel,
		request.BaseModelVersion,
		request.Date,
		request.NodeDID,
		request.SignedProof,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate aggregate model: %s", err), http.StatusInternalServerError)
		return
	}

	var response CalculationResponse
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse response: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ReadAggregatedModel handles retrieving a specific aggregated model
func ReadAggregatedModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := application_gateway.ReadAggregatedModel(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read aggregated model: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
}

// GetAllAggregatedModels handles retrieving all aggregated models
func GetAllAggregatedModels(w http.ResponseWriter, r *http.Request) {
	result, err := application_gateway.GetAllAggregatedModels()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get all aggregated models: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
}

// QueryAggregatedModelsByDateRange handles retrieving models by date range
func QueryAggregatedModelsByDateRange(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if startDate == "" || endDate == "" {
		http.Error(w, "Missing startDate or endDate parameters", http.StatusBadRequest)
		return
	}

	result, err := application_gateway.QueryAggregatedModelsByDateRange(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query aggregated models: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
} 