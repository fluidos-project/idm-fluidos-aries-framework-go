package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DHTDLTOperationsContract provides functions for DHT and DLT operations
type DHTDLTOperationsContract struct {
	contractapi.Contract
}

// AggregateModelTransaction represents a transaction in the blockchain (DLT)
type AggregateModelTransaction struct {
	ID              string        `json:"id"`
	Data            []interface{} `json:"data"`
	BaseModel       string        `json:"baseModel"`
	BaseModelVersion string       `json:"baseModelVersion"`
	Date            string        `json:"date"`
	NodeDID         string        `json:"nodeDID"`
	SignedProof     string        `json:"signedProof"`
	ModelsRef       []string      `json:"modelsRef"`
}

// Response struct for the calculation result
type CalculationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	TxID    string `json:"transactionId"`
	ID      string `json:"id"`
	ModelsRef []string `json:"modelsRef"`
}

// BaseModelUpdateTransaction represents a transaction in the blockchain (DLT) for base model updates
type BaseModelUpdateTransaction struct {
	ID              string   `json:"id"`
	BaseModel       string   `json:"baseModel"`
	BaseModelVersion string  `json:"baseModelVersion"`
	Date            string   `json:"date"`
	NodeDID         string   `json:"nodeDID"`
	SignedProof     string   `json:"signedProof"`
	DHTID           string   `json:"dhtId"`
}

type BaseModelUpdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	DHTID   string `json:"dhtId"`
}

// InitLedger adds a base set of modelUpdates to the ledger
func (s *DHTDLTOperationsContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	modelUpdates := []AggregateModelTransaction{}

	for _, modelUpdate := range modelUpdates {
		modelUpdateJSON, err := json.Marshal(modelUpdate)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(modelUpdate.ID, modelUpdateJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	baseModelUpdates := []BaseModelUpdateTransaction{}

	for _, baseModelUpdate := range baseModelUpdates {
		baseModelUpdateJSON, err := json.Marshal(baseModelUpdate)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(baseModelUpdate.ID, baseModelUpdateJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}


	return nil
}

// Retrieve information from the DHT and calculate the average model update and push it to the DLT
func (d *DHTDLTOperationsContract) AggregateModel(ctx contractapi.TransactionContextInterface, data []interface{}, baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string) (*CalculationResponse, error) {
	// Construct the ID
	id := fmt.Sprintf("%s_%s:%s", baseModel, baseModelVersion, date)

	exists, err := d.ModelUpdateExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("the model aggregation %s already exists", id)
	}

	// Validate data structure
	if data == nil {
		return nil, fmt.Errorf("data matrix cannot be nil")
	}

	// Mock modelsRef calculation (this will be replaced with actual DHT retrieval later)
	modelsRef := []string{
		fmt.Sprintf("model_%s_ref1", baseModel),
		fmt.Sprintf("model_%s_ref2", baseModel),
		fmt.Sprintf("model_%s_ref3", baseModel),
	}

	// Record the transaction in the DLT (blockchain)
	transaction := AggregateModelTransaction{
		ID:               id,
		Data:             data,
		BaseModel:        baseModel,
		BaseModelVersion: baseModelVersion,
		Date:            date,
		NodeDID:         nodeDID,
		SignedProof:     signedProof,
		ModelsRef:       modelsRef,
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(id, transactionJSON)
	if err != nil {
		return nil, err
	}

	txID := ctx.GetStub().GetTxID()

	response := &CalculationResponse{
		Status:    "success",
		Message:   fmt.Sprintf("Successfully stored aggregated model update with ID: %s", id),
		TxID:      txID,
		ID:        id,
		ModelsRef: modelsRef,
	}

	return response, nil
}

// ReadAggregatedModel returns the transaction stored in the world state with given id.
func (s *DHTDLTOperationsContract) ReadAggregatedModel(ctx contractapi.TransactionContextInterface, id string) (*AggregateModelTransaction, error) {
	transactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("the aggregated model %s does not exist", id)	
	}

	var transaction AggregateModelTransaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetAllAggregatedModels returns all transactions found in world state
func (s *DHTDLTOperationsContract) GetAllAggregatedModels(ctx contractapi.TransactionContextInterface) ([]*AggregateModelTransaction, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []*AggregateModelTransaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction AggregateModelTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// QueryAggregatedModelsByDateRange returns all transactions created between two timestamps
func (s *DHTDLTOperationsContract) QueryAggregatedModelsByDateRange(ctx contractapi.TransactionContextInterface, startDate string, endDate string) ([]*AggregateModelTransaction, error) {
	queryString := fmt.Sprintf(`{
		"selector": {
			"Timestamp": {
				"$gte": "%s",
				"$lte": "%s"
			}
		}
	}`, startDate, endDate)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer resultsIterator.Close()

	var transactions []*AggregateModelTransaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction AggregateModelTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// ModelUpdateExists returns true when transaction with given ID exists in world state
func (s *DHTDLTOperationsContract) ModelUpdateExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	modelUpdateJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return modelUpdateJSON != nil, nil
}

// CreateBaseModelUpdate adds a new base model update to the ledger
func (d *DHTDLTOperationsContract) CreateBaseModelUpdate(ctx contractapi.TransactionContextInterface, data []interface{}, baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string) (*BaseModelUpdateResponse, error) {
	// Add debug logging
	fmt.Printf("Creating base model update with parameters: baseModel=%s, version=%s, date=%s, nodeDID=%s\n", 
		baseModel, baseModelVersion, date, nodeDID)

	// Construct the ID
	id := fmt.Sprintf("%s_%s:%s:%s", baseModel, baseModelVersion, date, nodeDID)
	fmt.Printf("Generated ID: %s\n", id)

	// Check if transaction already exists
	exists, err := d.BaseModelUpdateExists(ctx, id)
	if err != nil {
		fmt.Printf("Error checking if base model exists: %v\n", err)
		return nil, fmt.Errorf("failed to check if base model exists: %v", err)
	}
	if exists {
		fmt.Printf("Base model update %s already exists\n", id)
		return nil, fmt.Errorf("base model update %s already exists", id)
	}

	// Mock DHT storage (to be implemented later)
	dhtID := fmt.Sprintf("dht_%s_%s", baseModel, date)
	fmt.Printf("Generated DHT ID: %s\n", dhtID)
	
	// Create the transaction
	transaction := BaseModelUpdateTransaction{
		ID:               id,
		BaseModel:        baseModel,
			BaseModelVersion: baseModelVersion,
			Date:            date,
			NodeDID:         nodeDID,
			SignedProof:     signedProof,
			DHTID:           dhtID,
	}

	// Marshal the transaction
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		fmt.Printf("Error marshaling transaction: %v\n", err)
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// Store in state
	err = ctx.GetStub().PutState(id, transactionJSON)
	if err != nil {
		fmt.Printf("Error putting state: %v\n", err)
		return nil, fmt.Errorf("failed to put state: %v", err)
	}

	response := &BaseModelUpdateResponse{
		Status:  "success",
		Message: fmt.Sprintf("Successfully stored base model update with ID: %s", id),
		DHTID:   dhtID,
	}

	fmt.Printf("Successfully created base model update: %+v\n", response)
	return response, nil
}

// BaseModelUpdateExists returns true when transaction with given ID exists in world state
func (s *DHTDLTOperationsContract) BaseModelUpdateExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	modelJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return modelJSON != nil, nil
}

// ReadBaseModelUpdate returns the transaction stored in the world state with given id.
func (s *DHTDLTOperationsContract) ReadBaseModelUpdate(ctx contractapi.TransactionContextInterface, id string) (*BaseModelUpdateTransaction, error) {
	modelJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if modelJSON == nil {
		return nil, fmt.Errorf("the base model update %s does not exist", id)
	}

	var transaction BaseModelUpdateTransaction
	err = json.Unmarshal(modelJSON, &transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetAllBaseModelUpdates returns all transactions found in world state
func (s *DHTDLTOperationsContract) GetAllBaseModelUpdates(ctx contractapi.TransactionContextInterface) ([]*BaseModelUpdateTransaction, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []*BaseModelUpdateTransaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction BaseModelUpdateTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			continue // Skip if not a base model update transaction
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// QueryBaseModelUpdatesByDateRange returns all transactions created between two timestamps
func (s *DHTDLTOperationsContract) QueryBaseModelUpdatesByDateRange(ctx contractapi.TransactionContextInterface, startDate string, endDate string) ([]*BaseModelUpdateTransaction, error) {
	queryString := fmt.Sprintf(`{
		"selector": {
			"Date": {
				"$gte": "%s",
				"$lte": "%s"
			}
		}
	}`, startDate, endDate)

	return s.queryBaseModelUpdates(ctx, queryString)
}

func (s *DHTDLTOperationsContract) queryBaseModelUpdates(ctx contractapi.TransactionContextInterface, queryString string) ([]*BaseModelUpdateTransaction, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []*BaseModelUpdateTransaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction BaseModelUpdateTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			continue
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func main() {
	modelTreatmentChaincode, err := contractapi.NewChaincode(&DHTDLTOperationsContract{})
	if err != nil {
		log.Panicf("Error creating model-treatment chaincode: %v", err)
	}

	if err := modelTreatmentChaincode.Start(); err != nil {
		log.Panicf("Error starting model-treatment chaincode: %v", err)
	}
}