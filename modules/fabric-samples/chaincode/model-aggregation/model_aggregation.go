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

// MockModelData represents a model from DHT with weights and ID
type MockModelData struct {
	Weights []interface{}
	DHTID   string
}

// calculateModelUpdate computes the difference between two matrices
func calculateModelUpdate(modelWeights, sourceWeights []interface{}) ([]float64, error) {
	var update []float64
	for i := 0; i < len(modelWeights); i++ {
		weight1, ok1 := modelWeights[i].(float64)
		weight2, ok2 := sourceWeights[i].(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("invalid matrix value type at index %d", i)
		}
		update = append(update, weight1-weight2)
	}
	return update, nil
}

// calculateAveragedUpdate computes the average of multiple model updates
func calculateAveragedUpdate(updates [][]float64) []float64 {
	if len(updates) == 0 {
		return nil
	}
	
	resultLen := len(updates[0])
	result := make([]float64, resultLen)
	
	for i := 0; i < resultLen; i++ {
		sum := 0.0
		for _, update := range updates {
			sum += update[i]
		}
		result[i] = sum / float64(len(updates))
	}
	
	return result
}

// calculateGlobalWeights computes final weights
func calculateGlobalWeights(averagedUpdate []float64, sourceWeights []interface{}) ([]float64, error) {
	var result []float64
	for i := 0; i < len(sourceWeights); i++ {
		sourceWeight, ok := sourceWeights[i].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid source weight type at index %d", i)
		}
		result = append(result, sourceWeight+averagedUpdate[i])
	}
	return result, nil
}

// generateMockWeights creates mock weights for testing
func generateMockWeights(length int) []interface{} {
	weights := make([]interface{}, length)
	for i := 0; i < length; i++ {
		weights[i] = float64(i) * 0.1 // Mock values
	}
	return weights
}

// generateMockModelData creates mock weights and DHTID for testing
func generateMockModelData(baseModel string, length int, index int) MockModelData {
	weights := make([]interface{}, length)
	for i := 0; i < length; i++ {
		weights[i] = float64(i) * 0.1 // Mock values
	}
	
	return MockModelData{
		Weights: weights,
		DHTID:   fmt.Sprintf("dht_%s_model_%d", baseModel, index),
	}
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

	// Mock DHT retrieval of other model weights with their IDs
	mockModels := []MockModelData{
		generateMockModelData(baseModel, len(data), 1),
		generateMockModelData(baseModel, len(data), 2),
		generateMockModelData(baseModel, len(data), 3),
	}

	// Calculate updates for each model
	var modelUpdates [][]float64
	var modelsRef []string // Store DHT IDs of participating models

	for _, model := range mockModels {
		update, err := calculateModelUpdate(model.Weights, data)
		if err != nil {
			return nil, fmt.Errorf("error calculating model update for model %s: %v", model.DHTID, err)
		}
		modelUpdates = append(modelUpdates, update)
		modelsRef = append(modelsRef, model.DHTID)
	}

	// Calculate averaged update
	averagedUpdate := calculateAveragedUpdate(modelUpdates)

	// Calculate global weights
	globalWeights, err := calculateGlobalWeights(averagedUpdate, data)
	if err != nil {
		return nil, fmt.Errorf("error calculating global weights: %v", err)
	}

	// Convert global weights to interface{} slice for storage
	var globalWeightsInterface []interface{}
	for _, w := range globalWeights {
		globalWeightsInterface = append(globalWeightsInterface, w)
	}

	// Create and store transaction
	transaction := AggregateModelTransaction{
		ID:              id,
		Data:            globalWeightsInterface,
		BaseModel:       baseModel,
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

	return &CalculationResponse{
		Status:    "success",
		Message:   fmt.Sprintf("Successfully aggregated model with ID: %s", id),
		TxID:      ctx.GetStub().GetTxID(),
		ID:        id,
		ModelsRef: modelsRef,
	}, nil
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