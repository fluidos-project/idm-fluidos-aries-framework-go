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

// DLTTransaction represents a transaction in the blockchain (DLT)
type AverageModelUpdateTransaction struct {
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

// InitLedger adds a base set of modelUpdates to the ledger
func (s *DHTDLTOperationsContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	modelUpdates := []AverageModelUpdateTransaction{}

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

	return nil
}

// Retrieve information from the DHT and calculate the average model update and push it to the DLT
func (d *DHTDLTOperationsContract) CalculateAverageModelUpdate(ctx contractapi.TransactionContextInterface, data []interface{}, baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string) (*CalculationResponse, error) {
	// Construct the ID
	id := fmt.Sprintf("%s_%s:%s", baseModel, baseModelVersion, date)

	exists, err := d.ModelUpdateExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("the model update %s already exists", id)
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
	transaction := AverageModelUpdateTransaction{
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
		Message:   fmt.Sprintf("Successfully stored model update with ID: %s", id),
		TxID:      txID,
		ID:        id,
		ModelsRef: modelsRef,
	}

	return response, nil
}

// ReadAverageModelUpdateTransaction returns the transaction stored in the world state with given id.
func (s *DHTDLTOperationsContract) ReadAverageModelUpdateTransaction(ctx contractapi.TransactionContextInterface, id string) (*AverageModelUpdateTransaction, error) {
	transactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("the transaction %s does not exist", id)	
	}

	var transaction AverageModelUpdateTransaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetAllAverageModelUpdateTransactions returns all transactions found in world state
func (s *DHTDLTOperationsContract) GetAllAverageModelUpdateTransactions(ctx contractapi.TransactionContextInterface) ([]*AverageModelUpdateTransaction, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []*AverageModelUpdateTransaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction AverageModelUpdateTransaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// QueryAverageModelUpdateTransactionsByDateRange returns all transactions created between two timestamps
func (s *DHTDLTOperationsContract) QueryAverageModelUpdateTransactionsByDateRange(ctx contractapi.TransactionContextInterface, startDate string, endDate string) ([]*AverageModelUpdateTransaction, error) {
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

	var transactions []*AverageModelUpdateTransaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction AverageModelUpdateTransaction
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

func main() {
	modelTreatmentChaincode, err := contractapi.NewChaincode(&DHTDLTOperationsContract{})
	if err != nil {
		log.Panicf("Error creating model-treatment chaincode: %v", err)
	}

	if err := modelTreatmentChaincode.Start(); err != nil {
		log.Panicf("Error starting model-treatment chaincode: %v", err)
	}
}