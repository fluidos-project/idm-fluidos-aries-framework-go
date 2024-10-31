package web

import (
	"encoding/json"
	application_gateway "fabricrest-go/application-gateway"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type AverageModelUpdateRequest struct {
	Data             []interface{} `json:"data"`
	BaseModel        string        `json:"baseModel"`
	BaseModelVersion string        `json:"baseModelVersion"`
	Date            string        `json:"date"`
	NodeDID         string        `json:"nodeDID"`
	SignedProof     string        `json:"signedProof"`
}

type CalculationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	TxID    string `json:"transactionId"`
	ID      string `json:"id"`
	ModelsRef []string `json:"modelsRef"`
}

// CalculateAverageModelUpdate handles the calculation and storage of average model updates
func CalculateAverageModelUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - CalculateAverageModelUpdate")

	var request AverageModelUpdateRequest
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

	result, err := application_gateway.CalculateAverageModelUpdate(
		request.Data,
		request.BaseModel,
		request.BaseModelVersion,
		request.Date,
		request.NodeDID,
		request.SignedProof,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate average model: %s", err), http.StatusInternalServerError)
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

// ReadAverageModelUpdate retrieves a specific model update transaction
func ReadAverageModelUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - ReadAverageModelUpdate")

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	transaction, err := application_gateway.ReadAverageModelUpdateTransaction(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read model update: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(transaction))
}

// GetAllModelUpdates retrieves all model update transactions
func GetAllModelUpdates(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - GetAllModelUpdates")

	transactions, err := application_gateway.GetAllAverageModelUpdateTransactions()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get all transactions: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(transactions))
}

// QueryModelUpdatesByDateRange retrieves transactions within a date range
func QueryModelUpdatesByDateRange(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - QueryModelUpdatesByDateRange")

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if startDate == "" || endDate == "" {
		http.Error(w, "Missing startDate or endDate parameters", http.StatusBadRequest)
		return
	}

	transactions, err := application_gateway.QueryAverageModelUpdateTransactionsByDateRange(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query transactions: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(transactions))
}
