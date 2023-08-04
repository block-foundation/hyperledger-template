package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestInit(t *testing.T) {
	scc := new(ItemContract)
	stub := shimtest.NewMockStub("item_contract_test", scc)

	res := stub.MockInit("000", nil)
	if res.Status != shim.OK {
		t.Errorf("Init failed: %s", res.Message)
	}
}

func TestInvoke(t *testing.T) {
	scc := new(ItemContract)
	stub := shimtest.NewMockStub("item_contract_test", scc)

	// Mock Init
	_ = stub.MockInit("000", nil)

	// Mock CreateItem
	item := &Item{ID: "item1", Name: "Laptop", Price: 1000}
	itemJSON, _ := json.Marshal(item)
	res := stub.MockInvoke("001", [][]byte{[]byte("createItem"), itemJSON})
	if res.Status != shim.OK {
		t.Errorf("Failed to create item: %s", res.Message)
	}

	// Mock GetItem
	res = stub.MockInvoke("002", [][]byte{[]byte("getItem"), []byte("item1")})
	if res.Status != shim.OK {
		t.Errorf("Failed to get item: %s", res.Message)
	}

	if !bytes.Equal(res.Payload, itemJSON) {
		t.Errorf("Returned item does not match expected item: %s != %s", string(res.Payload), string(itemJSON))
	}
}
