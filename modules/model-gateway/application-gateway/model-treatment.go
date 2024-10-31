package application_gateway

import (
	"encoding/json"
	"fmt"
)

// AggregateModel submits a transaction to aggregate model
func AggregateModel(data []interface{}, baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

	result, err := contract.SubmitTransaction(
		"AggregateModel",
		string(dataBytes),
		baseModel,
		baseModelVersion,
		date,
		nodeDID,
		signedProof,
	)
	if err != nil {
		return "", fmt.Errorf("failed to aggregate model: %w", err)
	}

	return string(result), nil
}

// ReadAggregatedModel retrieves a specific aggregated model
func ReadAggregatedModel(id string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")

	result, err := contract.EvaluateTransaction("ReadAggregatedModel", id)
	if err != nil {
		return "", fmt.Errorf("failed to read aggregated model: %w", err)
	}

	return string(result), nil
}

// GetAllAggregatedModels retrieves all aggregated models
func GetAllAggregatedModels() (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")

	result, err := contract.EvaluateTransaction("GetAllAggregatedModels")
	if err != nil {
		return "", fmt.Errorf("failed to get all aggregated models: %w", err)
	}

	return string(result), nil
}

// QueryAggregatedModelsByDateRange retrieves models within a date range
func QueryAggregatedModelsByDateRange(startDate string, endDate string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")

	result, err := contract.EvaluateTransaction("QueryAggregatedModelsByDateRange", startDate, endDate)
	if err != nil {
		return "", fmt.Errorf("failed to query aggregated models: %w", err)
	}

	return string(result), nil
}
