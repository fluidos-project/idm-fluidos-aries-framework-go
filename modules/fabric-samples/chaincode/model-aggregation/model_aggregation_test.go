package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "github.com/stretchr/testify/mock"
)

type MockContext struct {
    mock.Mock
    contractapi.TransactionContextInterface
}

func TestCalculateModelUpdate(t *testing.T) {
    // Test data
    modelWeights := [][]float64{{1.0, 2.0}, {3.0, 4.0}}
    sourceWeights := [][]float64{{0.5, 1.0}, {1.5, 2.0}}

    // Expected result
    expected := [][]float64{{0.5, 1.0}, {1.5, 2.0}}

    // Calculate update
    result, err := calculateModelUpdate(modelWeights, sourceWeights)

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}

func TestCalculateAveragedUpdate(t *testing.T) {
    // Test data
    updates := [][][]float64{
        {{1.0, 2.0}, {3.0, 4.0}},
        {{2.0, 3.0}, {4.0, 5.0}},
    }

    // Expected result
    expected := [][]float64{{1.5, 2.5}, {3.5, 4.5}}

    // Calculate average
    result := calculateAveragedUpdate(updates)

    // Assertions
    assert.Equal(t, expected, result)
}

func TestCalculateGlobalWeights(t *testing.T) {
    // Test data
    averagedUpdate := [][]float64{{0.5, 1.0}, {1.5, 2.0}}
    sourceWeights := [][]float64{{1.0, 2.0}, {3.0, 4.0}}

    // Expected result
    expected := [][]float64{{1.5, 3.0}, {4.5, 6.0}}

    // Calculate global weights
    result, err := calculateGlobalWeights(averagedUpdate, sourceWeights)

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}

func TestAggregateModel(t *testing.T) {
    contract := new(ModelAggregationContract)
    mockCtx := new(MockContext)

    // Test data
    baseModel := "testModel"
    baseModelVersion := "v1"
    date := "2024-03-01"
    nodeDID := "did:example:123"
    signedProof := "proof123"
    data := []interface{}{[]float64{1.0, 2.0}, []float64{3.0, 4.0}}

    // Call function
    response, err := contract.AggregateModel(mockCtx, data, baseModel, baseModelVersion, date, nodeDID, signedProof)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, "success", response.Status)
    assert.NotEmpty(t, response.ID)
    assert.NotEmpty(t, response.ModelsRef)
}

func TestQueryAggregatedModelsByDateRange(t *testing.T) {
    contract := new(ModelAggregationContract)
    mockCtx := new(MockContext)

    // Test data
    startDate := "2024-01-01"
    endDate := "2024-12-31"

    // Call function
    results, err := contract.QueryAggregatedModelsByDateRange(mockCtx, startDate, endDate)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, results)
}