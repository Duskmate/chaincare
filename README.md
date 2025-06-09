
# ChainCare - Blockchain Powered Healthcare on Hyperledger Fabric

**ChainCare** is a decentralized healthcare solution built on **Hyperledger Fabric**. It allows **hospitals** to securely store medical records and enables **insurance companies** to review, approve, or reject claims based on real-time, verifiable data, all while maintaining data integrity, privacy, and auditability.

## Features

- 🔐 Secure, tamper-proof medical records
- 🏥 Hospital: Record creation and updates medical records
- 🧾 Insurance Company: Claim validation and approval
- 🔄 Real time data access across organizations
- 🛡️ Full audit trail with Fabric ledger
- 📦 Simple scripted setup


## Installation (Linux Only)

Follow these steps to set up and interact with the ChainCare blockchain application.

### Step 1: Install Prerequisites & Clone Repo

```bash
  curl -s https://raw.githubusercontent.com/Duskmate/chaincare/refs/heads/main/install.sh | sudo bash
```
This script:

- Installs necessary dependencies (Fabric binaries, docker, etc)
- Clones the chaincare repository
- Sets up your environment for Fabric development

### Step 2: Navigate to the Project Directory

```bash
  cd chaincare
```
### Step 3: Start the Network and Create the Channel

```bash
  sudo ./network.sh up createChannel -c mednet
```
This spins up your Fabric network, creates the mednet channel, and connects all organizations (Hospital, Insurance) into the blockchain channel.

### Step 4: Deploy the Chaincode (Smart Contract)

```bash
  sudo ./network.sh deployCC -c mednet -ccn medblock -ccp ./chaincode -ccl go
```

Deploys the medcc chaincode written in Go, enabling secure transactions for managing medical records and claims.
    
## Try It Out | Interacting with the Network

### Step 5: Set Context to Hospital Org (Org1)

```bash
  source ./scripts/setOrgPeerContext.sh 1
```

### Step 6: Add a Medical Record

```bash
  peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C mednet -n medblock \
  --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_HOSPITAL_CA \
  --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_INSURANCE_CA \
  -c '{"function": "AddMedicalRecord", "Args":["rec001", "patient1", "Qm123abcXYZ"]}'
```

### Step 7: Retrieve a Medical Record

```bash
  peer chaincode query \
  -C mednet -n medblock \
  -c '{"function": "GetMedicalRecord", "Args":["rec001"]}'
```

### Step 8: Set Context to Insurance Org (Org2)

```bash
  source ./scripts/setOrgPeerContext.sh 2
```

### Step 9: Request an Insurance Claim

```bash
  peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C mednet -n medblock \
  --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_HOSPITAL_CA \
  --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_INSURANCE_CA \
  -c '{"function": "RequestClaim", "Args":["claim001", "rec001"]}'
```

### Step 10: Approve the Claim

```bash
  peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C mednet -n medblock \
  --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_HOSPITAL_CA \
  --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_INSURANCE_CA \
  -c '{"function": "ApproveClaim", "Args":["claim001"]}'
```

### Step 11: Reject a Claim

```bash
  peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C mednet -n medblock \
  --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_HOSPITAL_CA \
  --peerAddresses localhost:9051 --tlsRootCertFiles $PEER0_INSURANCE_CA \
  -c '{"function": "RejectClaim", "Args":["claim001"]}'
```

### Step 12: View Claim Status

```bash
  peer chaincode query \
  -C mednet -n medblock \
  -c '{"function": "GetClaim", "Args":["claim001"]}'
```
## Chaincode Functions (Quick Overview)

| Function         | Org              | Description                       |
| ---------------- | ---------------- | --------------------------------- |
| AddMedicalRecord | Hospital (Org1)  | Adds new encrypted patient record |
| GetMedicalRecord | Any              | View medical record info          |
| RequestClaim     | Insurance (Org2) | Create claim request for a record |
| ApproveClaim     | Insurance (Org2) | Approves a claim after validation |
| RejectClaim      | Insurance (Org2) | Rejects a claim                   |
| GetClaim         | Any              | Query claim status                |


## Tech Stack

- 🔗 Hyperledger Fabric v2.5
- 🌐 Go (Chaincode)
- 🐳 Docker / Docker Compose
- 🐧 Linux Shell Automation
- 🧠 IPFS-ready Design (optional for record hashes)


## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/)
## Authors

- [@duskmate](https://github.com/Duskmate/)
