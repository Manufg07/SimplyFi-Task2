# 🚚 Supply Chain Tracking Chaincode

A Hyperledger Fabric-based solution for tracking products through supply chain stages with enforced state transitions and complete audit history.

## ✨ Features

- 📦 **Product Lifecycle Management**: Track products from manufacture to delivery
- 🛠️ **State Transition Validation**: Enforce valid status changes between stages
- 🕰️ **Immutable History**: Maintain timestamped record of all status changes
- 🔍 **Advanced Queries**: Filter by status, retrieve full history, and list all products
- 🛡️ **Data Integrity**: Comprehensive input validation and error handling

## ⚙️ Prerequisites

- 🐳 Hyperledger Fabric v2.4+ network
- 🐹 Go 1.23 environment
- 🧰 Fabric SDK dependencies (`fabric-contract-api-go/v2`)
- 🔐 Properly configured cryptographic materials


# 💡Fabric Test Network commands

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
# Install/update dependencies in the chaincode folder
```
go mod tidy
```

## 🚢 Deploy chaincode listed in the samples.
```
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-supplychain -ccl go
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
## 1. Register a new product
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
-C mychannel -n basic \
--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
-c '{"function":"RegisterProduct","Args":["PROD-001"]}'
```
## 2. Verify registration
```
peer chaincode query -C mychannel -n basic -c '{"function":"GetProduct","Args":["PROD-001"]}'
```
## 3. Update to Shipped status
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
-C mychannel -n basic \
--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
-c '{"function":"UpdateStatus","Args":["PROD-001","Shipped"]}'
```

## 4. Check current status
```
peer chaincode query -C mychannel -n basic -c '{"function":"GetProduct","Args":["PROD-001"]}'
```

## 5. Update to In-Transit status
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
-C mychannel -n basic \
--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
-c '{"function":"UpdateStatus","Args":["PROD-001","In-Transit"]}'
```

## 6. View status history
```
peer chaincode query -C mychannel -n basic -c '{"function":"GetProductHistory","Args":["PROD-001"]}'
```

## 7. Mark as Delivered
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
-C mychannel -n basic \
--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
-c '{"function":"UpdateStatus","Args":["PROD-001","Delivered"]}'
```

## 8. List all products
```
peer chaincode query -C mychannel -n basic -c '{"function":"GetAllProducts","Args":[]}'
```

## 9. Filter by Delivered status
```
peer chaincode query -C mychannel -n basic -c '{"function":"GetProductsByStatus","Args":["Delivered"]}'
```

## 10. Test invalid transition (should fail)
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
-C mychannel -n basic \
--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
-c '{"function":"UpdateStatus","Args":["PROD-001","Manufactured"]}'
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

## 🚀 Performance Optimization -CouchDB indexes(Future Enhancement) or We can use Composite key method

For production deployments or large-scale implementations, consider adding CouchDB indexing 
to optimize query performance:

## create directory
```bash
mkdir -p chaincode-supplychain/META-INF/statedb/couchdb/indexes
```
## statusIndex.json file
```
{
    "index": {
        "fields": ["currentStatus"]
    },
    "name": "statusIndex",
    "type": "json"
}
```
### 📁 CouchDB Index Setup
1. **Enable CouchDB** in your network:
```bash
./network.sh up createChannel -ca -s couchdb

