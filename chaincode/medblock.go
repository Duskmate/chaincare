package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type MedicalRecord struct {
	ID        string `json:"id"`
	PatientID string `json:"patientId"`
	IPFSHash  string `json:"ipfsHash"`
	Timestamp string `json:"timestamp"`
	CreatedBy string `json:"createdBy"`
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

// Utility: Get client MSP ID
func getClientMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	return ctx.GetClientIdentity().GetMSPID()
}

// Utility: Get client identity's common name
func getClientID(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", err
	}
	parts := strings.Split(id, "/CN=")
	if len(parts) > 1 {
		return parts[1], nil
	}
	return id, nil
}

// === Chaincode Functions ===

func (s *SmartContract) AddMedicalRecord(ctx contractapi.TransactionContextInterface, recordID, patientID, ipfsHash string) error {
	mspID, _ := getClientMSPID(ctx)
	if !strings.Contains(strings.ToLower(mspID), "hospital") {
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

func (s *SmartContract) GetMedicalRecord(ctx contractapi.TransactionContextInterface, recordID string) (*MedicalRecord, error) {
	data, err := ctx.GetStub().GetState("RECORD_" + recordID)
	if err != nil || data == nil {
		return nil, fmt.Errorf("record not found")
	}
	var record MedicalRecord
	err = json.Unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *SmartContract) RequestClaim(ctx contractapi.TransactionContextInterface, claimID, recordID string) error {
	mspID, _ := getClientMSPID(ctx)
	if !(strings.Contains(strings.ToLower(mspID), "patient") || strings.Contains(strings.ToLower(mspID), "insurance")) {
		return fmt.Errorf("only patients or insurance can request claims")
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
	mspID, _ := getClientMSPID(ctx)
	if !strings.Contains(strings.ToLower(mspID), "insurance") {
		return fmt.Errorf("only insurance can approve claims")
	}

	claimBytes, err := ctx.GetStub().GetState("CLAIM_" + claimID)
	if err != nil || claimBytes == nil {
		return fmt.Errorf("claim not found")
	}

	var claim Claim
	json.Unmarshal(claimBytes, &claim)

	if claim.Status == Rejected {
		return fmt.Errorf("claims once rejected cannot be approved. please initiate a new claim")
	}

	claim.Status = Approved
	claim.ApprovedBy, _ = getClientID(ctx)

	updatedBytes, _ := json.Marshal(claim)
	return ctx.GetStub().PutState("CLAIM_"+claimID, updatedBytes)
}

func (s *SmartContract) RejectClaim(ctx contractapi.TransactionContextInterface, claimID string) error {
	mspID, _ := getClientMSPID(ctx)
	if !strings.Contains(strings.ToLower(mspID), "insurance") {
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

func (s *SmartContract) GetClaim(ctx contractapi.TransactionContextInterface, claimID string) (*Claim, error) {
	claimBytes, err := ctx.GetStub().GetState("CLAIM_" + claimID)
	if err != nil || claimBytes == nil {
		return nil, fmt.Errorf("claim not found")
	}

	var claim Claim
	err = json.Unmarshal(claimBytes, &claim)
	if err != nil {
		return nil, err
	}
	return &claim, nil
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
