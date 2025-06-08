package main

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type MedicalRecord struct {
	ID         string `json:"id"`
	PatientID  string `json:"patientId"`
	IPFSHash   string `json:"ipfsHash"`
	Timestamp  string `json:"timestamp"`
	CreatedBy  string `json:"createdBy"`
}

type ClaimStatus string

const (
	Pending  ClaimStatus = "PENDING"
	Approved ClaimStatus = "APPROVED"
	Rejected ClaimStatus = "REJECTED"
)

type Claim struct {
	ClaimID    string      `json:"claimId"`
	RecordID   string      `json:"recordId"`
	Status     ClaimStatus `json:"status"`
	ApprovedBy string      `json:"approvedBy"`
}

// === Utility Functions ===

// Get client identity's full ID
func getClientID(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", err
	}
	return id, nil
}

// Get client MSP ID
func getClientMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	return mspid, nil
}

// === Chaincode Functions ===

func (s *SmartContract) AddMedicalRecord(ctx contractapi.TransactionContextInterface, recordID, patientID, ipfsHash string) error {
	mspid, err := getClientMSPID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MSP ID: %v", err)
	}
	if mspid != "HospitalMSP" {
		return fmt.Errorf("only hospitals can add records")
	}

	creator, _ := getClientID(ctx)
	record := MedicalRecord{
		ID:        recordID,
		PatientID: patientID,
		IPFSHash:  ipfsHash,
		Timestamp: time.Now().Format(time.RFC3339),
		CreatedBy: creator,
	}

	recordBytes, _ := json.Marshal(record)
	return ctx.GetStub().PutState("RECORD_"+recordID, recordBytes)
}

func (s *SmartContract) ViewMedicalRecords(ctx contractapi.TransactionContextInterface) ([]*MedicalRecord, error) {
	clientID, _ := getClientID(ctx)
	query := fmt.Sprintf(`{"selector":{"patientId":"%s"}}`, clientID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*MedicalRecord
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()
		var record MedicalRecord
		json.Unmarshal(queryResponse.Value, &record)
		records = append(records, &record)
	}
	return records, nil
}

func (s *SmartContract) RequestClaim(ctx contractapi.TransactionContextInterface, claimID, recordID string) error {
	mspid, err := getClientMSPID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MSP ID: %v", err)
	}
	if mspid != "PatientMSP" {
		return fmt.Errorf("only patients can request claims")
	}

	claim := Claim{
		ClaimID:    claimID,
		RecordID:   recordID,
		Status:     Pending,
		ApprovedBy: "",
	}

	claimBytes, _ := json.Marshal(claim)
	return ctx.GetStub().PutState("CLAIM_"+claimID, claimBytes)
}

func (s *SmartContract) ApproveClaim(ctx contractapi.TransactionContextInterface, claimID string) error {
	mspid, err := getClientMSPID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MSP ID: %v", err)
	}
	if mspid != "InsuranceMSP" {
		return fmt.Errorf("only insurance can approve claims")
	}

	claimBytes, err := ctx.GetStub().GetState("CLAIM_" + claimID)
	if err != nil || claimBytes == nil {
		return fmt.Errorf("claim not found")
	}

	var claim Claim
	json.Unmarshal(claimBytes, &claim)

	claim.Status = Approved
	claim.ApprovedBy, _ = getClientID(ctx)

	updatedBytes, _ := json.Marshal(claim)
	return ctx.GetStub().PutState("CLAIM_"+claimID, updatedBytes)
}

func (s *SmartContract) RejectClaim(ctx contractapi.TransactionContextInterface, claimID string) error {
	mspid, err := getClientMSPID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MSP ID: %v", err)
	}
	if mspid != "InsuranceMSP" {
		return fmt.Errorf("only insurance can reject claims")
	}

	claimBytes, err := ctx.GetStub().GetState("CLAIM_" + claimID)
	if err != nil || claimBytes == nil {
		return fmt.Errorf("claim not found")
	}

	var claim Claim
	json.Unmarshal(claimBytes, &claim)

	claim.Status = Rejected
	claim.ApprovedBy, _ = getClientID(ctx)

	updatedBytes, _ := json.Marshal(claim)
	return ctx.GetStub().PutState("CLAIM_"+claimID, updatedBytes)
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		panic(fmt.Sprintf("Error creating MedBlock chaincode: %s", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Failed to start chaincode: %s", err))
	}
}