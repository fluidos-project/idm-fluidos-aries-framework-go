package application_gateway

import (
	"fmt"
)



// DHT FEA


func writeToDHT() {
	if err := InitializeConnection(); err != nil {
		return false, fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()
	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("model-treatment")
	


}

func readFromDHT(key string) (string, error) {
    if err := InitializeConnection(); err != nil {
        return "", fmt.Errorf("failed to initialize blockchain connection: %v", err)
    }
    defer CloseConnection()

    network := gateway.GetNetwork("mychannel")
    contract := network.GetContract("model-treatment")

    // Call the ReadFromDHT function in the smart contract
    result, err := contract.EvaluateTransaction("ReadFromDHT", key)
    if err != nil {
        return "", fmt.Errorf("failed to read from DHT: %w", err)
    }

    return string(result), nil
}

// FLUIDOS
func SetAuthReq(timestamp string, action string, resource string, id string, subject string, decision string) (bool, error) {
	if err := InitializeConnection(); err != nil {
		return false, fmt.Errorf("failed to initialize blockchain connection: %v", err)
	}
	defer CloseConnection()

	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("fluidosAccessHist")

	_, err := contract.SubmitTransaction("CreateAsset", timestamp, action, resource, id, subject, decision)
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
