package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"githubcom/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStub implements the ChaincodeStubInterface for testing
type MockStub struct {
	mock.Mock
	shim.ChaincodeStubInterface
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)
	return args.Error(0)
}

// MockContext implements the TransactionContextInterface for testing
type MockContext struct {
	mock.Mock
	contractapi.TransactionContextInterface
	stub *MockStub
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	return mc.stub
}

func TestCreateBaseModelUpdate(t *testing.T) {
	// Create a new instance of our contract
	contract := new(DHTDLTOperationsContract)
	
	// Create our mock stub and context
	mockStub := new(MockStub)
	mockContext := new(MockContext)
	mockContext.stub = mockStub

	// Test data
	baseModel := "testModel"
	baseModelVersion := "v1"
	date := "2024-03-01"
	nodeDID := "did:example:123"
	signedProof := "proof123"
	data := []interface{}{[]float64{1.0, 2.0}, []float64{3.0, 4.0}}

	// Expected ID
	expectedID := "testModel_v1:2024-03-01:did:example:123"

	// Setup mock expectations
	mockStub.On("GetState", expectedID).Return([]byte{}, nil)
	mockStub.On("PutState", expectedID, mock.Anything).Return(nil)

	// Call the function we want to test
	response, err := contract.CreateBaseModelUpdate(mockContext, data, baseModel, baseModelVersion, date, nodeDID, signedProof)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "success", response.Status)
	assert.Contains(t, response.Message, expectedID)
	assert.NotEmpty(t, response.DHTID)

	// Verify mock expectations
	mockStub.AssertExpectations(t)
}

func TestReadBaseModelUpdate(t *testing.T) {
	contract := new(DHTDLTOperationsContract)
	mockStub := new(MockStub)
	mockContext := new(MockContext)
	mockContext.stub = mockStub

	// Test data
	id := "testModel_v1:2024-03-01:did:example:123"
	expectedTransaction := &BaseModelUpdateTransaction{
		ID:              id,
		BaseModel:       "testModel",
		BaseModelVersion: "v1",
		Date:            "2024-03-01",
		NodeDID:         "did:example:123",
		SignedProof:     "proof123",
		DHTID:           "dht_testModel_2024-03-01",
	}

	// Marshal the expected transaction
	transactionJSON, _ := json.Marshal(expectedTransaction)

	// Setup mock expectations
	mockStub.On("GetState", id).Return(transactionJSON, nil)

	// Call the function
	result, err := contract.ReadBaseModelUpdate(mockContext, id)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTransaction.ID, result.ID)
	assert.Equal(t, expectedTransaction.BaseModel, result.BaseModel)
	assert.Equal(t, expectedTransaction.DHTID, result.DHTID)

	// Verify mock expectations
	mockStub.AssertExpectations(t)
}

func TestInitLedger(t *testing.T) {
	contract := new(DHTDLTOperationsContract)
	mockStub := new(MockStub)
	mockContext := new(MockContext)
	mockContext.stub = mockStub

	// Setup mock expectations for each PutState call
	mockStub.On("PutState", mock.Anything, mock.Anything).Return(nil)

	// Call InitLedger
	err := contract.InitLedger(mockContext)

	// Assertions
	assert.NoError(t, err)

	// Verify that PutState was called the expected number of times
	// (depends on how many items you initialize in InitLedger)
	mockStub.AssertNumberOfCalls(t, "PutState", 2) // Adjust this number based on your InitLedger implementation
} 