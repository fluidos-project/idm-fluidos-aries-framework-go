package application_gateway

import (
	"fmt"
)

// FLUIDOS
func SetAuthReq(timestamp string, action string, resource string, id string, did string, subject string, decision string) (bool, error) {
	if err := InitializeConnection(); err != nil {
		return false, fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("fluidosAccessHist")

	_, err := contract.SubmitTransaction("CreateAsset", timestamp, action, resource, id, did, subject, decision)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return true, nil
}

func GetAuthReq(key string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("fluidosAccessHist")

	// Evaluate transaction
	authReq, err := contract.EvaluateTransaction("ReadAsset", key)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return string(authReq), nil

}

func GetAuthReqsByDate(startDate string, endDate string) (string, error) {
	if err := InitializeConnection(); err != nil {
		return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("fluidosAccessHist")

	// Evaluate transaction
	authReqs, err := contract.EvaluateTransaction("QueryAssetsByDateRange", startDate, endDate)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return string(authReqs), nil
}
