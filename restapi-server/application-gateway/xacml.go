package application_gateway

import (
	"fmt"
	"strings"
)

// XACML
func SetXACMLEntity(key string, value string) (bool, error) {
	if err := InitializeConnection(); err != nil {
		return false, fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("xacml")

	_, err := contract.SubmitTransaction("Set", key, value)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return true, nil
}

func GetXACMLEntity(key string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("xacml")

	// Evaluate transaction
	entity, err := contract.EvaluateTransaction("Get", key)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return string(entity), nil
}

func IndexQueryXACML(query string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("xacml")

	// Evaluate transaction
	queryRes, err := contract.EvaluateTransaction("QueryAssets", query)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	queryResString := string(queryRes)
	// Replace
	queryResString = strings.ReplaceAll(queryResString, "\\\"", "\"")
	queryResString = strings.ReplaceAll(queryResString, "\\u003e", ">")
	queryResString = strings.ReplaceAll(queryResString, "\\u003c", "<")
	queryResString = strings.ReplaceAll(queryResString, "\\\\", "\\")
	queryResString = strings.ReplaceAll(queryResString, "\"{", "{")
	queryResString = strings.ReplaceAll(queryResString, "}\"", "}")
	queryResString = strings.ReplaceAll(queryResString, "\\<", "<")
	queryResString = strings.ReplaceAll(queryResString, "\\>", ">")

	return queryResString, nil
}

func UpdateEntity(key string, value string) (bool, error) {
	if err := InitializeConnection(); err != nil {
		return false, fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("xacml")

	_, err := contract.SubmitTransaction("Update", key, value)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return true, nil
}
