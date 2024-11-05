package web

import (
	"encoding/json"
	application_gateway "fabricrest-go/application-gateway"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type BaseModelUpdateRequest struct {
	BaseModel        string        `json:"baseModel"`
	BaseModelVersion string        `json:"baseModelVersion"`
	Date            string        `json:"date"`
	NodeDID         string        `json:"nodeDID"`
	SignedProof     string        `json:"signedProof"`
	Data            []interface{} `json:"data"`
}

type BaseModelUpdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	DHTID   string `json:"dhtId"`
}

// CreateBaseModelUpdate handles creating a new base model update
func CreateBaseModelUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - CreateBaseModelUpdate")

	var request BaseModelUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if request.Data == nil {
		fmt.Println("Error: Data matrix is nil")
		http.Error(w, "Data matrix cannot be nil", http.StatusBadRequest)
		return
	}
	if request.BaseModel == "" {
		fmt.Println("Error: BaseModel is empty")
		http.Error(w, "BaseModel is required", http.StatusBadRequest)
		return
	}
	if request.BaseModelVersion == "" {
		fmt.Println("Error: BaseModelVersion is empty")
		http.Error(w, "BaseModelVersion is required", http.StatusBadRequest)
		return
	}

	fmt.Println("Validation passed, calling gateway")

	result, err := application_gateway.CreateBaseModelUpdate(
		request.BaseModel,
		request.BaseModelVersion,
		request.Date,
		request.NodeDID,
		request.SignedProof,
		request.Data,
	)
	if err != nil {
		fmt.Printf("Error from gateway: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to create base model update: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Gateway returned successfully: %s\n", result)

	var response BaseModelUpdateResponse
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to parse response: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ReadBaseModelUpdate handles retrieving a specific base model update
func ReadBaseModelUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := application_gateway.ReadBaseModelUpdate(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read base model update: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
}

// GetAllBaseModelUpdates handles retrieving all base model updates
func GetAllBaseModelUpdates(w http.ResponseWriter, r *http.Request) {
	result, err := application_gateway.GetAllBaseModelUpdates()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get all base model updates: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
} 

// QueryBaseModelUpdatesByDateRange handles retrieving base model updates by date range
func QueryBaseModelUpdatesByDateRange(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if startDate == "" || endDate == "" {
		http.Error(w, "Missing startDate or endDate parameters", http.StatusBadRequest)
		return
	}

	result, err := application_gateway.QueryBaseModelUpdatesByDateRange(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query base model updates: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(result))
}