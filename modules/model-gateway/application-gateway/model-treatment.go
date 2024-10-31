package application_gateway

import (
	"encoding/json"
	"fmt"
)

// CalculateAverageModelUpdate submits a transaction to calculate and store average model update
func CalculateAverageModelUpdate(data []interface{}, baseModel string, baseModelVersion string, date string, nodeDID string, signedProof string) (string, error) {
    if err := InitializeConnection(); err != nil {
        return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
    }
    defer CloseConnection()

    network := gateway.GetNetwork("mychannel")
    contract := network.GetContract("model-treatment")

    // Convert data to JSON string for chaincode
    dataBytes, err := json.Marshal(data)
    if err != nil {
        return "", fmt.Errorf("failed to marshal data: %v", err)
    }

    result, err := contract.SubmitTransaction(
        "CalculateAverageModelUpdate",
        string(dataBytes),
        baseModel,
        baseModelVersion,
        date,
        nodeDID,
        signedProof,
    )
    if err != nil {
        return "", fmt.Errorf("failed to calculate average model update: %w", err)
    }

    return string(result), nil
}

// ReadAverageModelUpdateTransaction retrieves a specific model update transaction
func ReadAverageModelUpdateTransaction(id string) (string, error) {
    if err := InitializeConnection(); err != nil {
        return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
    }
    defer CloseConnection()

    network := gateway.GetNetwork("mychannel")
    contract := network.GetContract("model-treatment")

    result, err := contract.EvaluateTransaction("ReadAverageModelUpdateTransaction", id)
    if err != nil {
        return "", fmt.Errorf("failed to read model update transaction: %w", err)
    }

    return string(result), nil
}

// GetAllAverageModelUpdateTransactions retrieves all model update transactions
func GetAllAverageModelUpdateTransactions() (string, error) {
    if err := InitializeConnection(); err != nil {
        return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
    }
    defer CloseConnection()

    network := gateway.GetNetwork("mychannel")
    contract := network.GetContract("model-treatment")

    result, err := contract.EvaluateTransaction("GetAllAverageModelUpdateTransactions")
    if err != nil {
        return "", fmt.Errorf("failed to get all model update transactions: %w", err)
    }

    return string(result), nil
}

// QueryAverageModelUpdateTransactionsByDateRange retrieves transactions within a date range
func QueryAverageModelUpdateTransactionsByDateRange(startDate string, endDate string) (string, error) {
    if err := InitializeConnection(); err != nil {
        return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
    }
    defer CloseConnection()

    network := gateway.GetNetwork("mychannel")
    contract := network.GetContract("model-treatment")

    result, err := contract.EvaluateTransaction("QueryAverageModelUpdateTransactionsByDateRange", startDate, endDate)
    if err != nil {
        return "", fmt.Errorf("failed to query model update transactions: %w", err)
    }

    return string(result), nil
}
