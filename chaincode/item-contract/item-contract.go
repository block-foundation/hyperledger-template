package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Item. It includes a reference
// to contractapi.Contract to inherit standard Smart Contract functionality including
// functions like TransactionContextInterface.
type SmartContract struct {
	contractapi.Contract
}

// Item describes basic details of what makes up a simple item.
// This is the data that will be stored in the ledger for each item.
type Item struct {
	ID    string `json:"id"`    // The ID is the key in the ledger
	Name  string `json:"name"`  // Name will be used to describe the item
	Price int    `json:"price"` // Price will hold the cost of each item
}

// InitLedger adds a base set of items to the ledger.
// This function is called when the Smart Contract is instantiated.
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	items := []Item{
		{ID: "item1", Name: "Item 1", Price: 100},
		{ID: "item2", Name: "Item 2", Price: 200},
		{ID: "item3", Name: "Item 3", Price: 300},
	}

	// Add each item to the ledger
	for _, item := range items {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return err
		}

		// The PutState function writes data to the ledger, using the item ID as key.
		err = ctx.GetStub().PutState(item.ID, itemJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateItem issues a new item to the world state with given details.
// This function can be called by a client application to add new items to the ledger.
func (s *SmartContract) CreateItem(ctx contractapi.TransactionContextInterface, id string, name string, price int) error {
	item := Item{
		ID:    id,
		Name:  name,
		Price: price,
	}

	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	// The item is added to the ledger with its ID as the key.
	return ctx.GetStub().PutState(id, itemJSON)
}

// ReadItem returns the item stored in the world state with given id.
// It reads the item from the ledger using the GetState function.
func (s *SmartContract) ReadItem(ctx contractapi.TransactionContextInterface, id string) (*Item, error) {
	itemJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %v", err)
	}

	if itemJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var item Item
	err = json.Unmarshal(itemJSON, &item)
	if err != nil {
		return nil, err
	}

	// Returns the item read from the ledger.
	return &item, nil
}

// UpdateItem updates an existing item in the world state with provided parameters.
// This function can be called by a client application to modify an item in the ledger.
func (s *SmartContract) UpdateItem(ctx contractapi.TransactionContextInterface, id string, name string, price int) error {
	item, err := s.ReadItem(ctx, id)
	if err != nil {
		return err
	}

	item.Name = name
	item.Price = price

	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err
	}

	// The updated item is written back to the ledger.
	return ctx.GetStub().PutState(id, itemJSON)
}

// DeleteItem deletes an given item from the world state.
// This function can be called by a client application to remove an item from the ledger.
func (s *SmartContract) DeleteItem(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.ItemExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// The DelState function is used to remove the item from the ledger.
	return ctx.GetStub().DelState(id)
}

// ItemExists returns true when item with given ID exists in world state.
// This function is used inside the DeleteItem function to check if an item exists before attempting to delete it.
func (s *SmartContract) ItemExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	itemJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state. %v", err)
	}

	return itemJSON != nil, nil
}

// GetAllItems returns all items found in the world state.
// This function can be used by a client application to list all items in the ledger.
func (s *SmartContract) GetAllItems(ctx contractapi.TransactionContextInterface) ([]*Item, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var items []*Item
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var item Item
		err = json.Unmarshal(queryResponse.Value, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	// Returns all items found.
	return items, nil
}

// main function starts up the chaincode in the container during instantiate.
func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create example chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting example chaincode: %s", err.Error())
	}
}
