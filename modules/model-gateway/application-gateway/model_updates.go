package application_gateway

import (
	"encoding/json"
	"fmt"
)

func CreateBaseModelUpdate(baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string, data []interface{}) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-updates")

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

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
		return "", fmt.Errorf("failed to create base model update: %w", err)
	}

	return string(result), nil
}

func ReadBaseModelUpdate(id string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-updates")

	result, err := contract.EvaluateTransaction("ReadBaseModelUpdate", id)
	if err != nil {
		return "", fmt.Errorf("failed to read base model update: %w", err)
	}

	return string(result), nil
}

func GetAllBaseModelUpdates() (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-updates")

	result, err := contract.EvaluateTransaction("GetAllBaseModelUpdates")
	if err != nil {
		return "", fmt.Errorf("failed to get all base model updates: %w", err)
	}

	return string(result), nil
}

func QueryBaseModelUpdatesByDateRange(startDate string, endDate string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-updates")

	result, err := contract.EvaluateTransaction("QueryBaseModelUpdatesByDateRange", startDate, endDate)
	if err != nil {
		return "", fmt.Errorf("failed to query base model updates: %w", err)
	}

	return string(result), nil
} 