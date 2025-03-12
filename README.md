# 🚀 SimplyFi Task2

- Basic Chaincode Implementation for Asset Management
- Chaincode for Supply Chain Tracking

![Hyperledger Fabric](https://img.shields.io/badge/Hyperledger-Fabric-2F3134?logo=hyperledger&logoColor=white)

## 📚 Table of Contents
- [Hyperledger Fabric Setup](#-hyperledger-fabric-setup)
- [Test Network Deployment](#-test-network-deployment)
- [Chaincode Interaction](#-chaincode-interaction)
- [Cleanup](#-cleanup)

# 🔧 Hyperledger Fabric Setup

Reference: https://hyperledger-fabric.readthedocs.io/en/latest/index.html

### Commands to setup fabric-samples for running a test network and deplying an chaincode. Other required Commands are given in the curresponding folders readme.
- Step1 : setup fabric samples
- Step2 : clone the repo and add chaincode-Asset folder inside CHF/fabric-samples/asset-transfer-basic/ chaincode-supplychain
-Step 3: Follow the commands inside the Chaincode-Asset Readme 
-Step 4: Same for chaincode-supplychain

## ⚙️ Install Dependencies

Note: If any of the following dependencies are available on your laptop, then no need to install it.

## cURL
Install curl using the command
```
sudo apt install curl -y
```

To verify the installation enter the following command

```
curl -V
```

## JQ
Install JQ using the following command
```
sudo apt install jq -y
```

To verify the installation enter the following command


```
jq --version
```

## Build Essential
Install Build Essential uisng the commnad
```
sudo apt install build-essential
```
To verify the installation enter the following command


```
dpkg -l | grep build-essential

```

### Download Fabric 

Note: Open a terminal in the **CHF** Folder & Execute the Following Commands

`curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh`

`./install-fabric.sh -f '2.5.4' -c '1.5.7'`

`sudo cp fabric-samples/bin/* /usr/local/bin`
