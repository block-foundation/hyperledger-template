<div align="right">

[![GitHub License](https://img.shields.io/github/license/block-foundation/blocktxt?style=flat-square&logo=readthedocs&logoColor=FFFFFF&label=&labelColor=%23041B26&color=%23041B26&link=LICENSE)](https://github.com/block-foundation/hyperledger-template/blob/main/LICENSE)
[![devContainer](https://img.shields.io/badge/Container-Remote?style=flat-square&logo=visualstudiocode&logoColor=%23FFFFFF&label=Remote&labelColor=%23041B26&color=%23041B26)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/block-foundation/hyperledger-template)

</div>

---

<div>
    <img align="right" src="https://raw.githubusercontent.com/block-foundation/brand/master/src/logo/logo_gray.png" width="96" alt="Block Foundation Logo">
    <h1 align="left">Hyperledger Template</h1>
    <h3 align="left">Block Foundation</h3>
</div>

---

<img align="right" width="75%" src="https://raw.githubusercontent.com/block-foundation/brand/master/src/image/repository_cover/block_foundation-structure-03-accent.jpg"  alt="Block Foundation Brand">

### Contents

- [Introduction](#introduction)
- [Colophon](#colophon)

<br clear="both"/>

---

<div align="right">

[![Report a Bug](https://img.shields.io/badge/Report%20a%20Bug-GitHub?style=flat-square&&logoColor=%23FFFFFF&color=%23E1E4E5)](https://github.com/block-foundation/hyperledger-template/issues/new?assignees=&labels=Needs%3A+Triage+%3Amag%3A%2Ctype%3Abug-suspected&projects=&template=bug_report.yml)
[![Request a Feature](https://img.shields.io/badge/Request%20a%20Feature-GitHub?style=flat-square&&logoColor=%23FFFFFF&color=%23E1E4E5)](https://github.com/block-foundation/hyperledger-template/issues/new?assignees=&labels=Needs%3A+Triage+%3Amag%3A%2Ctype%3Abug-suspected&projects=&template=feature_request.yml)
[![Ask a Question](https://img.shields.io/badge/Ask%20a%20Question-GitHub?style=flat-square&&logoColor=%23FFFFFF&color=%23E1E4E5)](https://github.com/block-foundation/hyperledger-template/issues/new?assignees=&labels=Needs%3A+Triage+%3Amag%3A%2Ctype%3Abug-suspected&projects=&template=question.yml)
[![Make a Suggestion](https://img.shields.io/badge/Make%20a%20Suggestion-GitHub?style=flat-square&&logoColor=%23FFFFFF&color=%23E1E4E5)](https://github.com/block-foundation/hyperledger-template/issues/new?assignees=&labels=Needs%3A+Triage+%3Amag%3A%2Ctype%3Abug-suspected&projects=&template=suggestion.yml)
[![Start a Discussion](https://img.shields.io/badge/Start%20a%20Discussion-GitHub?style=flat-square&&logoColor=%23FFFFFF&color=%23E1E4E5)](https://github.com/block-foundation/hyperledger-template/issues/new?assignees=&labels=Needs%3A+Triage+%3Amag%3A%2Ctype%3Abug-suspected&projects=&template=discussion.yml)

</div>

## Introduction

Below is a basic example of a chaincode (smart contract) in Hyperledger Fabric. This is written in Go, which is one of the languages you can use to create chaincode in Hyperledger Fabric. This simple example chaincode models assets (items) that can be created, read, updated, deleted, and listed.

``` go
package main

import (
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

type Item struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Price int    `json:"price"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    items := []Item{
        {ID: "item1", Name: "Item 1", Price: 100},
        {ID: "item2", Name: "Item 2", Price: 200},
        {ID: "item3", Name: "Item 3", Price: 300},
    }

    for _, item := range items {
        itemJSON, err := json.Marshal(item)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(item.ID, itemJSON)
        if err != nil {
            return fmt.Errorf("failed to put to world state. %v", err)
        }
    }

    return nil
}

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

    return ctx.GetStub().PutState(id, itemJSON)
}

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

    return &item, nil
}

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

    return ctx.GetStub().PutState(id, itemJSON)
}

func (s *SmartContract) DeleteItem(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.ItemExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
}

func (s *SmartContract) ItemExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    itemJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state. %v", err)
    }

    return itemJSON != nil, nil
}

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

    return items, nil
}

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
```

This chaincode allows for the following operations:

1. Initialize the ledger with some items (`InitLedger` function).
2. Create an item (`CreateItem` function).
3. Read an item (`ReadItem` function).
4. Check if an item exists (`ItemExists` function).
5. Update an item (`UpdateItem` function).
6. Delete an item (`DeleteItem` function).
7. Get all items (`GetAllItems` function).

Each item has an ID, a name, and a price. Remember to use the correct import paths for the `fmt`, `json`, and `github.com/hyperledger/fabric-contract-api-go/contractapi` packages.



## Directory structure

The directory structure for a Hyperledger Fabric Chaincode (smart contract):

``` sh
- fabric-item-contract
    - chaincode
        - item-contract
            - item-contract.go // The chaincode file
    - test
        - item-contract_test.go // Unit tests for your smart contract
    - scripts
        - startFabric.sh // Scripts for setting up and running the network
        - testAPI.sh     // Scripts for testing the API
    - doc
        - item-contract.md // Documentation for your chaincode
    - network
        - // Contains the configuration for your network
    - client
        - // Optional, for client side code if needed
    - README.md // Main project readme
    - .gitignore // For excluding files from version control
```

This is a fairly standard structure for a Fabric chaincode project.

- Your actual chaincode would live inside the `chaincode` directory.
- The `test` directory contains the tests for your chaincode.
- The `scripts` directory contains shell scripts for setting up your environment and running your application.
- The `doc` directory contains the documentation for your chaincode.
- The `network` directory would contain the configurations for your network.
- The `client` directory is optional and would contain the client side code for invoking the chaincode, if any.

These directories and their contents will help keep your project organized and will make it easier for others to understand your work.






---

## Colophon

### Authors

This is an open-source project by the **[Block Foundation](https://www.blockfoundation.io "Block Foundation website")**.

The Block Foundation mission is enabling architects to take back initiative and contribute in solving the mismatch in housing through blockchain technology. Therefore the Block Foundation seeks to unschackle the traditional constraints and construct middle ground between rent and the rigidity of traditional mortgages.

website: [www.blockfoundation.io](https://www.blockfoundation.io "Block Foundation website")

### Development Resources

#### Contributing

We'd love for you to contribute and to make this project even better than it is today!
Please refer to the [contribution guidelines](.github/CONTRIBUTING.md) for information.

### Legal Information

#### Copyright

Copyright &copy; 2023 [Stichting Block Foundation](https://www.blockfoundation.io/ "Block Foundation website"). All Rights Reserved.

#### License

Except as otherwise noted, the content in this repository is licensed under the
[Creative Commons Attribution 4.0 International (CC BY 4.0) License](https://creativecommons.org/licenses/by/4.0/), and
code samples are licensed under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).

Also see [LICENSE](https://github.com/block-foundation/community/blob/master/src/LICENSE) and [LICENSE-CODE](https://github.com/block-foundation/community/blob/master/src/LICENSE-CODE).

#### Disclaimer

**THIS SOFTWARE IS PROVIDED AS IS WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING ANY IMPLIED WARRANTIES OF FITNESS FOR A PARTICULAR PURPOSE, MERCHANTABILITY, OR NON-INFRINGEMENT.**
