package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ModelAggregationContract provides functions for model aggregation operations
type ModelAggregationContract struct {
	contractapi.Contract
}

// AggregateModelTransaction represents a transaction in the blockchain
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

// CalculationResponse represents the response for model aggregation
type CalculationResponse struct {
	Status    string   `json:"status"`
	Message   string   `json:"message"`
	TxID      string   `json:"transactionId"`
	ID        string   `json:"id"`
	ModelsRef []string `json:"modelsRef"`
}

// InitLedger adds a base set of model aggregations to the ledger
func (s *ModelAggregationContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	modelUpdates := []AggregateModelTransaction{
		{
			ID:               "basemodel1_v1:2024-03-01",
			BaseModel:        "basemodel1",
			BaseModelVersion: "v1",
			Date:            "2024-03-01",
			NodeDID:         "did:example:node1",
			SignedProof:     "proof1",
			ModelsRef:       []string{"model_basemodel1_ref1", "model_basemodel1_ref2", "model_basemodel1_ref3"},
			Data:            []interface{}{},
		},
	}

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

// AggregateModel creates a new aggregated model in the ledger
func (s *ModelAggregationContract) AggregateModel(ctx contractapi.TransactionContextInterface, data []interface{}, baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string) (*CalculationResponse, error) {
	id := fmt.Sprintf("%s_%s:%s", baseModel, baseModelVersion, date)

	exists, err := s.ModelExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("the model aggregation %s already exists", id)
	}

	if data == nil {
		return nil, fmt.Errorf("data matrix cannot be nil")
	}

	modelsRef := []string{
		fmt.Sprintf("model_%s_ref1", baseModel),
		fmt.Sprintf("model_%s_ref2", baseModel),
		fmt.Sprintf("model_%s_ref3", baseModel),
	}

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
		Message:   fmt.Sprintf("Successfully stored aggregated model with ID: %s", id),
		TxID:      txID,
		ID:        id,
		ModelsRef: modelsRef,
	}

	return response, nil
}

// ReadAggregatedModel retrieves a specific model by ID
func (s *ModelAggregationContract) ReadAggregatedModel(ctx contractapi.TransactionContextInterface, id string) (*AggregateModelTransaction, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("the model %s does not exist", id)
	}

	var transaction AggregateModelTransaction
	err = json.Unmarshal(data, &transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// ModelExists returns true when model with given ID exists
func (s *ModelAggregationContract) ModelExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	modelJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return modelJSON != nil, nil
}

// GetAllAggregatedModels returns all models
func (s *ModelAggregationContract) GetAllAggregatedModels(ctx contractapi.TransactionContextInterface) ([]*AggregateModelTransaction, error) {
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
					continue
			}
			transactions = append(transactions, &transaction)
	}

	return transactions, nil
}


// QueryModelsByDateRange returns models within the specified date range
func (s *ModelAggregationContract) QueryAggregatedModelsByDateRange(ctx contractapi.TransactionContextInterface, startDate string, endDate string) ([]*AggregateModelTransaction, error) {
	queryString := fmt.Sprintf(`{
		"selector": {
			"Date": {
				"$gte": "%s",
				"$lte": "%s"
			}
		}
	}`, startDate, endDate)

	return s.queryModels(ctx, queryString)
}

func (s *ModelAggregationContract) queryModels(ctx contractapi.TransactionContextInterface, queryString string) ([]*AggregateModelTransaction, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
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
			continue
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}



func main() {
	modelAggregationChaincode, err := contractapi.NewChaincode(&ModelAggregationContract{})
	if err != nil {
		log.Panicf("Error creating model-aggregation chaincode: %v", err)
	}

	if err := modelAggregationChaincode.Start(); err != nil {
		log.Panicf("Error starting model-aggregation chaincode: %v", err)
	}
} 