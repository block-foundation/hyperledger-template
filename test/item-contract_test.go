// This is the test file for the chaincode defined in item-contract.go. The functions in this file
// use the Go 'testing' package to define test cases for the smart contract functions.

package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
)

// TestInitLedger tests the InitLedger function of the smart contract. It first initializes the chaincode
// and the mock stub. Then it invokes the InitLedger function and checks the response status to be 'OK'.
// After that, it verifies that the "item1" is created correctly in the ledger.
func TestInitLedger(t *testing.T) {
	// Initialize the chaincode and stub
	chaincode := new(SmartContract)
	stub := shimtest.NewMockStub("test", chaincode)

	// Invoke InitLedger
	res := stub.MockInvoke("tx1", [][]byte{[]byte("InitLedger")})

	// Check the response
	assert.Equal(t, int32(shim.OK), res.Status, "InitLedger failed")

	// Check if item1 was created
	itemBytes := stub.State["item1"]
	if itemBytes == nil {
		t.Fatal("Failed to create item1 in InitLedger")
	}

	item := new(Item)
	_ = json.Unmarshal(itemBytes, item)

	// Verify if the item was created with the correct details
	assert.Equal(t, item, &Item{ID: "item1", Name: "Item 1", Price: 100})
}

// TestCreateItem tests the CreateItem function of the smart contract. It first initializes the chaincode
// and the mock stub. Then it invokes the CreateItem function and checks the response status to be 'OK'.
// After that, it verifies that the "item4" is created correctly in the ledger.
func TestCreateItem(t *testing.T) {
	// Initialize the chaincode and stub
	chaincode := new(SmartContract)
	stub := shimtest.NewMockStub("test", chaincode)

	// Invoke CreateItem
	res := stub.MockInvoke("tx1", [][]byte{[]byte("CreateItem"), []byte("item4"), []byte("Item 4"), []byte("400")})

	// Check the response
	assert.Equal(t, int32(shim.OK), res.Status, "CreateItem failed")

	// Check if item4 was created
	itemBytes := stub.State["item4"]
	if itemBytes == nil {
		t.Fatal("Failed to create item4 in CreateItem")
	}

	item := new(Item)
	_ = json.Unmarshal(itemBytes, item)

	// Verify if the item was created with the correct details
	assert.Equal(t, item, &Item{ID: "item4", Name: "Item 4", Price: 400})
}

// TestReadItem tests the ReadItem function of the smart contract. It first initializes the chaincode
// and the mock stub. Then it invokes the ReadItem function and checks the response status to be 'OK'.
// After that, it verifies that the "item1" is read correctly from the ledger.
func TestReadItem(t *testing.T) {
	// Initialize the chaincode and stub
	chaincode := new(SmartContract)
	stub := shimtest.NewMockStub("test", chaincode)

	// Invoke ReadItem
	res := stub.MockInvoke("tx1", [][]byte{[]byte("ReadItem"), []byte("item1")})

	// Check the response
	assert.Equal(t, int32(shim.OK), res.Status, "ReadItem failed")

	item := new(Item)
	_ = json.Unmarshal(res.Payload, item)

	// Verify if the item read is correct
	assert.Equal(t, item, &Item{ID: "item1", Name: "Item 1", Price: 100})
}
