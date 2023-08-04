#!/bin/bash
#
# This script starts a basic Hyperledger Fabric network consisting of a single organization
# with one peer and a single-channel.

# This script does the following:

# 1. Sets the shell to exit if any command fails and to print each command before execution.
# 2. Starts the Hyperledger Fabric network.
# 3. Creates and joins a channel.
# 4. Launches the CLI container.
# 5. Installs and instantiates the chaincode.

# Note: The actual startFabric.sh script would depend on your specific use case, environment, and network configuration. Always customize and test your scripts to ensure they work as expected in your environment.

# This script assumes that you are using the fabric-samples provided by the Hyperledger Fabric project, have a network.sh script available to spin up the network and have a docker-compose.yml for your containers. Also, it assumes your chaincode name is "item-contract" and you've put your chaincode in github.com/chaincode/item-contract/go. Make sure to replace these with your actual chaincode name and path.


# Exit on first error, print all commands.
set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

starttime=$(date +%s)

# Launch network; create a channel and join peer to channel
cd ../network
./network.sh up createChannel -c mychannel -ca

# Now launch the CLI container in order to install, instantiate chaincode
docker-compose -f ./docker-compose.yml up -d cli

docker exec cli peer chaincode install -n item-contract -v 0 -p github.com/chaincode/item-contract/go
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n item-contract -v 0 -c '{"Args":[]}' -P "AND('Org1MSP.peer')"

printf "\nTotal execution time : $(($(date +%s) - starttime)) secs ...\n\n"
echo "Hyperledger Fabric network started"
