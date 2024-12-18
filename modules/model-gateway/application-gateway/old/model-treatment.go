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

// CreateBaseModelUpdate submits a transaction to create a base model update
func CreateBaseModelUpdate(baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string, data []interface{}) (string, error) {
	fmt.Printf("Starting CreateBaseModelUpdate in gateway\n")

	if err := InitializeConnection(); err != nil {
		fmt.Printf("Failed to initialize connection: %v\n", err)
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	fmt.Printf("Connection initialized successfully\n")

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")

	fmt.Printf("Got contract reference\n")

	// Convert data to JSON string
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to marshal data: %v\n", err)
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

	fmt.Printf("Submitting transaction with parameters: baseModel=%s, version=%s, date=%s, nodeDID=%s\n",
		baseModel, baseModelVersion, date, nodeDID)

	result, err := contract.SubmitTransaction(
		"CreateBaseModelUpdate",
		string(dataBytes),
		baseModel,
		baseModelVersion,
		date,
		nodeDID,
		signedProof,
	)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %v\n", err)
		return "", fmt.Errorf("failed to create base model update: %w", err)
	}

	fmt.Printf("Transaction submitted successfully. Result: %s\n", string(result))
	return string(result), nil
}

// ReadBaseModelUpdate retrieves a specific base model update
func ReadBaseModelUpdate(id string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")

	result, err := contract.EvaluateTransaction("ReadBaseModelUpdate", id)
	if err != nil {
		return "", fmt.Errorf("failed to read base model update: %w", err)
	}

	return string(result), nil
}

// GetAllBaseModelUpdates retrieves all base model updates
func GetAllBaseModelUpdates() (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")

	result, err := contract.EvaluateTransaction("GetAllBaseModelUpdates")
	if err != nil {
		return "", fmt.Errorf("failed to get all base model updates: %w", err)
	}

	return string(result), nil
}

// QueryBaseModelUpdatesByDateRange retrieves base model updates within a date range
func QueryBaseModelUpdatesByDateRange(startDate string, endDate string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")

	result, err := contract.EvaluateTransaction("QueryBaseModelUpdatesByDateRange", startDate, endDate)
	if err != nil {
		return "", fmt.Errorf("failed to query base model updates: %w", err)
	}

	return string(result), nil
}
