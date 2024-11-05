package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ModelUpdatesContract provides functions for base model update operations
type ModelUpdatesContract struct {
	contractapi.Contract
}

// BaseModelUpdateTransaction represents a transaction in the blockchain
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

// InitLedger adds a base set of model updates to the ledger
func (s *ModelUpdatesContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	baseModelUpdates := []BaseModelUpdateTransaction{
		{
			ID:               "basemodel1_v1:2024-03-01:did:example:node1",
			BaseModel:        "basemodel1",
			BaseModelVersion: "v1",
			Date:            "2024-03-01",
			NodeDID:         "did:example:node1",
			SignedProof:     "proof1",
			DHTID:           "dht_basemodel1_2024-03-01",
		},
	}

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

// CreateBaseModelUpdate creates a new base model update in the ledger
func (s *ModelUpdatesContract) CreateBaseModelUpdate(ctx contractapi.TransactionContextInterface, data []interface{}, baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string) (*BaseModelUpdateResponse, error) {
	id := fmt.Sprintf("%s_%s:%s:%s", baseModel, baseModelVersion, date, nodeDID)

	exists, err := s.BaseModelUpdateExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("base model update %s already exists", id)
	}

	dhtID := fmt.Sprintf("dht_%s_%s", baseModel, date)

	transaction := BaseModelUpdateTransaction{
		ID:               id,
		BaseModel:        baseModel,
		BaseModelVersion: baseModelVersion,
		Date:            date,
		NodeDID:         nodeDID,
		SignedProof:     signedProof,
		DHTID:           dhtID,
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	err = ctx.GetStub().PutState(id, transactionJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to put state: %v", err)
	}

	response := &BaseModelUpdateResponse{
		Status:  "success",
		Message: fmt.Sprintf("Successfully stored base model update with ID: %s", id),
		DHTID:   dhtID,
	}

	return response, nil
}

// ReadBaseModelUpdate returns the base model update stored in the world state
func (s *ModelUpdatesContract) ReadBaseModelUpdate(ctx contractapi.TransactionContextInterface, id string) (*BaseModelUpdateTransaction, error) {
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

// BaseModelUpdateExists returns true when model with given ID exists
func (s *ModelUpdatesContract) BaseModelUpdateExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	modelJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return modelJSON != nil, nil
}

// GetAllBaseModelUpdates returns all base model updates
func (s *ModelUpdatesContract) GetAllBaseModelUpdates(ctx contractapi.TransactionContextInterface) ([]*BaseModelUpdateTransaction, error) {
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
			continue
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// QueryBaseModelUpdatesByDateRange returns models within the specified date range
func (s *ModelUpdatesContract) QueryBaseModelUpdatesByDateRange(ctx contractapi.TransactionContextInterface, startDate string, endDate string) ([]*BaseModelUpdateTransaction, error) {
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

func (s *ModelUpdatesContract) queryBaseModelUpdates(ctx contractapi.TransactionContextInterface, queryString string) ([]*BaseModelUpdateTransaction, error) {
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
	modelUpdatesChaincode, err := contractapi.NewChaincode(&ModelUpdatesContract{})
	if err != nil {
		log.Panicf("Error creating model-updates chaincode: %v", err)
	}

	if err := modelUpdatesChaincode.Start(); err != nil {
		log.Panicf("Error starting model-updates chaincode: %v", err)
	}
} 