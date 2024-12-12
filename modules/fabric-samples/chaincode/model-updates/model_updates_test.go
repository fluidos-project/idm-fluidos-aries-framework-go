package main

import (
    "testing"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockStub implements the ChaincodeStubInterface for testing
type MockStub struct {
    mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
    args := ms.Called(key)
    return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
    args := ms.Called(key, value)
    return args.Error(0)
}

func (ms *MockStub) GetStateByRange(startKey, endKey string) (contractapi.StateQueryIteratorInterface, error) {
    args := ms.Called(startKey, endKey)
    return args.Get(0).(contractapi.StateQueryIteratorInterface), args.Error(1)
}

func (ms *MockStub) GetQueryResult(query string) (contractapi.StateQueryIteratorInterface, error) {
    args := ms.Called(query)
    return args.Get(0).(contractapi.StateQueryIteratorInterface), args.Error(1)
}

// MockContext implements the TransactionContextInterface for testing
type MockContext struct {
    mock.Mock
    contractapi.TransactionContextInterface
    stub *MockStub
}

func (mc *MockContext) GetStub() *MockStub {
    return mc.stub
}

func TestInitLedger(t *testing.T) {
    contract := new(ModelUpdatesContract)
    mockStub := new(MockStub)
    mockContext := &MockContext{stub: mockStub}

    // Setup expectations
    mockStub.On("PutState", mock.Anything, mock.Anything).Return(nil)

    // Test InitLedger
    err := contract.InitLedger(mockContext)

    // Assertions
    assert.NoError(t, err)
    mockStub.AssertExpectations(t)
}

func TestCreateBaseModelUpdate(t *testing.T) {
    contract := new(ModelUpdatesContract)
    mockStub := new(MockStub)
    mockContext := &MockContext{stub: mockStub}

    // Test data
    baseModel := "testModel"
    baseModelVersion := "v1"
    date := "2024-03-01"
    nodeDID := "did:example:123"
    signedProof := "proof123"
    data := []interface{}{[]float64{1.0, 2.0}, []float64{3.0, 4.0}}

    // Expected ID
    expectedID := "testModel_v1:2024-03-01:did:example:123"

    // Setup expectations
    mockStub.On("GetState", expectedID).Return([]byte{}, nil)
    mockStub.On("PutState", expectedID, mock.Anything).Return(nil)

    // Test CreateBaseModelUpdate
    response, err := contract.CreateBaseModelUpdate(mockContext, data, baseModel, baseModelVersion, date, nodeDID, signedProof)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, "success", response.Status)
    assert.Contains(t, response.Message, expectedID)
    assert.NotEmpty(t, response.DHTID)
    mockStub.AssertExpectations(t)
}

func TestReadBaseModelUpdate(t *testing.T) {
    contract := new(ModelUpdatesContract)
    mockStub := new(MockStub)
    mockContext := &MockContext{stub: mockStub}

    // Test data
    id := "testModel_v1:2024-03-01:did:example:123"
    expectedTransaction := &BaseModelUpdateTransaction{
        ID:               id,
        BaseModel:        "testModel",
        BaseModelVersion: "v1",
        Date:            "2024-03-01",
        NodeDID:         "did:example:123",
        SignedProof:     "proof123",
        DHTID:           "dht_testModel_2024-03-01",
    }

    // Marshal the expected transaction
    transactionJSON, _ := json.Marshal(expectedTransaction)

    // Setup expectations
    mockStub.On("GetState", id).Return(transactionJSON, nil)

    // Test ReadBaseModelUpdate
    result, err := contract.ReadBaseModelUpdate(mockContext, id)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedTransaction.ID, result.ID)
    assert.Equal(t, expectedTransaction.BaseModel, result.BaseModel)
    assert.Equal(t, expectedTransaction.DHTID, result.DHTID)
    mockStub.AssertExpectations(t)
}

func TestQueryBaseModelUpdatesByDateRange(t *testing.T) {
    contract := new(ModelUpdatesContract)
    mockStub := new(MockStub)
    mockContext := &MockContext{stub: mockStub}

    // Test data
    startDate := "2024-01-01"
    endDate := "2024-12-31"
    expectedQuery := `{
        "selector": {
            "Date": {
                "$gte": "2024-01-01",
                "$lte": "2024-12-31"
            }
        }
    }`

    // Setup mock iterator
    mockIterator := new(MockQueryIterator)
    mockStub.On("GetQueryResult", mock.Anything).Return(mockIterator, nil)
    mockIterator.On("HasNext").Return(false)
    mockIterator.On("Close").Return(nil)

    // Test QueryBaseModelUpdatesByDateRange
    results, err := contract.QueryBaseModelUpdatesByDateRange(mockContext, startDate, endDate)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, results)
    mockStub.AssertExpectations(t)
    mockIterator.AssertExpectations(t)
}

// MockQueryIterator implements StateQueryIteratorInterface for testing
type MockQueryIterator struct {
    mock.Mock
    contractapi.StateQueryIteratorInterface
}

func (mqi *MockQueryIterator) HasNext() bool {
    args := mqi.Called()
    return args.Bool(0)
}

func (mqi *MockQueryIterator) Next() (*contractapi.QueryResponse, error) {
    args := mqi.Called()
    return args.Get(0).(*contractapi.QueryResponse), args.Error(1)
}

func (mqi *MockQueryIterator) Close() error {
    args := mqi.Called()
    return args.Error(0)
}