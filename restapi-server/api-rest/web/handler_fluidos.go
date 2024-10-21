package web

import (
	"encoding/json"
	application_gateway "fabricrest-go/application-gateway"
	"fmt"
	"io"
	"net/http"
)

type AuthReqPayload struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Subject   string `json:"subject"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Decision  string `json:"decision"`
}

type QueryByDateRange struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// FLUIDOS
func RegisterAuthReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - RegisterAuthReq")
	var payload AuthReqPayload

	//Verify if HTTP method is valid
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	//Validate the required fields
	if payload.Id == "" || payload.Timestamp == "" || payload.Subject == "" ||
		payload.Action == "" || payload.Resource == "" || payload.Decision == "" {
		http.Error(w, "Missing parameters in JSON", http.StatusBadRequest)
		return
	}

	registered, err := application_gateway.SetAuthReq(payload.Timestamp, payload.Action, payload.Resource, payload.Id, payload.Subject, payload.Decision)

	if registered {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Authorization Request '%s' successfully registered in Hyperledger Fabric\n", payload.Id)
	} else {
		fmt.Fprintf(w, "Failed to register Authorization Request '%s' in Fabric: %s", payload.Id, err)
		http.Error(w, fmt.Sprintf("There was a problem while registering the Authorization Request '%s' in Fabric", payload.Id), http.StatusInternalServerError)
		return
	}
}

func QueryAuthReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - QueryAuthReq")

	//Verify if HTTP method is valid
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}

	authReq, err := application_gateway.GetAuthReq(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authorization Request with id '%s' not found", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(authReq))
}

func QueryAuthReqByDate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request - QueryAuthReqsByDate")
	var payload QueryByDateRange

	//Verify if HTTP method is valid
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if payload.StartDate == "" || payload.EndDate == "" {
		http.Error(w, "Missing parameters in JSON", http.StatusBadRequest)
		return
	}

	authReqs, err := application_gateway.GetAuthReqsByDate(payload.StartDate, payload.EndDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("There was an error in Blockchain: %s", err), http.StatusInternalServerError)
		return
	}

	if authReqs == "" {
		http.Error(w, fmt.Sprintf("No Authorization Requests found between '%s' and '%s'", payload.StartDate, payload.EndDate), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(authReqs))
}
