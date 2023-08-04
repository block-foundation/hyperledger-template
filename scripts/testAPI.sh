#!/bin/bash

# This script is used to test the API endpoints for the Hyperledger Fabric item contract.
# It assumes you have a server running on localhost:8080 that exposes the necessary endpoints.

# This script does the following:

# 1. Sets the shell to exit if any command fails and to print each command before execution.
# 2. Makes a POST request to create a new item.
# 3. Makes a GET request to read the created item.
# 4. Makes a PUT request to update the created item.
# 5. Makes a DELETE request to delete the updated item.
# 6. Makes a GET request to list all items.

# Exit on first error, print all commands.
set -ev

# Test API to create a new item
curl -X POST http://localhost:8080/items -H 'Content-Type: application/json' -d '{
    "id": "item5",
    "name": "Item 5",
    "price": 500
}'

# Test API to get an item
curl -X GET http://localhost:8080/items/item5

# Test API to update an item
curl -X PUT http://localhost:8080/items/item5 -H 'Content-Type: application/json' -d '{
    "name": "Updated Item 5",
    "price": 550
}'

# Test API to delete an item
curl -X DELETE http://localhost:8080/items/item5

# Test API to get all items
curl -X GET http://localhost:8080/items

