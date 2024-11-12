package application_gateway

import (
	"encoding/json"
	"fmt"
)

func AggregateModel(data [][]float64, baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string) (string, error) {
	// Convert [][]float64 to []interface{} for chaincode
	dataInterface := make([]interface{}, len(data))
	for i, row := range data {
		dataInterface[i] = row
	}

	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-aggregation")

	dataBytes, err := json.Marshal(dataInterface)
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

func ReadAggregatedModel(id string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-aggregation")

	result, err := contract.EvaluateTransaction("ReadAggregatedModel", id)
	if err != nil {
		return "", fmt.Errorf("failed to read aggregated model: %w", err)
	}

	return string(result), nil
}

func GetAllAggregatedModels() (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-aggregation")

	result, err := contract.EvaluateTransaction("GetAllAggregatedModels")
	if err != nil {
		return "", fmt.Errorf("failed to get all aggregated models: %w", err)
	}

	return string(result), nil
}

func QueryAggregatedModelsByDateRange(startDate string, endDate string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-aggregation")

	result, err := contract.EvaluateTransaction("QueryAggregatedModelsByDateRange", startDate, endDate)
	if err != nil {
		return "", fmt.Errorf("failed to query aggregated models: %w", err)
	}

	return string(result), nil
} 