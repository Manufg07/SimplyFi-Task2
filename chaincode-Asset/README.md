# 🏦 Asset Management Chaincode

A Hyperledger Fabric-based chaincode for managing digital assets on a blockchain network. Provides full CRUD operations and owner-based query capabilities.

## ✨ Features

- ✅ **Create Asset**: Register new assets on the ledger
- 🔍 **Read Asset**: Retrieve asset details by ID
- ✏️ **Update Asset**: Modify existing asset properties
- 🗑️ **Delete Asset**: Remove assets from the ledger
- 👤 **Query by Owner**: Find all assets owned by a specific entity
- 📢 **Event System**: Emits blockchain events for critical operations
- 🛡️ **Data Validation**: Comprehensive input validation and error handling

## 🚀 Core Functions

### 📜 Main Chaincode Functions
| Function Name            | Description                                                                 |
|--------------------------|-----------------------------------------------------------------------------|
| `CreateAsset`            | Registers new assets on the ledger with initial properties                  |
| `ReadAsset`              | Retrieves asset details by unique ID                                        |
| `UpdateAsset`            | Modifies existing asset properties (name, owner, value)                     |
| `DeleteAsset`            | Permanently removes an asset from the ledger                                |                                      |
| `QueryAssetsByOwner`     | Finds all assets owned by a specific entity                                 |
| `AssetExists`            | Checks existence of an asset in the ledger                                  |

## ⚙️ Prerequisites

- 🐳 Hyperledger Fabric v2.4+ network
- 🐹 Go 1.23.0 environment
- 🧰 Fabric SDK dependencies
- 🔐 Properly configured cryptographic materials

# 💡 Fabric Test Network commands

Note: Open a terminal in the **CHF** Folder & Execute the Following Commands
```
cd fabric-samples/test-network/
```

./network.sh -h

## Bring up the network using CA

```
./network.sh up createChannel -ca

```
```
docker ps -a
```

## 🚢 Deploy  chaincode listed in the samples.
```
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-Asset -ccl go
```
# 🧪 Testing

## Open a new terminal and Setup the environment variables
```
export FABRIC_CFG_PATH=$PWD/../config/

export CORE_PEER_TLS_ENABLED=true

export CORE_PEER_LOCALMSPID="Org1MSP"

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

export CORE_PEER_ADDRESS=localhost:7051
```
## 1.Invoke command to add assets
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"CreateAsset","Args":["Asset-1","Laptop","user1","1000"]}'
```
## 2.Query Command to read an asset
```
peer chaincode query -C mychannel -n basic -c '{"function":"ReadAsset","Args":["Asset-1"]}'

```
## 3.Invoke command to Update assets
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"UpdateAsset","Args":["Asset-1","Updated Laptop","user1","1200"]}'
```
## 4.Invoke command to Delete assets
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"DeleteAsset","Args":["Asset-1"]}'
```
## 5.Query Command to check Asset exist
```
peer chaincode query -C mychannel -n basic -c '{"function":"AssetExists","Args":["Asset-1"]}'

```

## 6.Query Command to read QueryAssetsByOwner
```
peer chaincode query -C mychannel -n basic -c '{"function":"QueryAssetsByOwner","Args":["user1"]}'
```


# 🧹 Clean Up

## Tear down network
```
./network.sh down
```
### Remove chaincode containers
```
docker rm -f $(docker ps -aq)
```

### Remove chaincode images
```
docker rmi -f $(docker images -q)
```


