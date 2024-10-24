package contract

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, transactionJSON string) error {
	var transaction models.Transaction
	err := json.Unmarshal([]byte(transactionJSON), &transaction)
	if err != nil {
		return fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	transactionDid := transaction.Did.ToString()
	readTransaction, _ := s.ReadTransaction(ctx, transactionDid)

	if readTransaction != nil {
		return fmt.Errorf("the transaction %s already exists", transactionDid)
	}

	return ctx.GetStub().PutState(transactionDid, []byte(transactionJSON))
}

func (s *SmartContract) ReadTransaction(ctx contractapi.TransactionContextInterface, did string) (*models.Transaction, error) {
	transactionJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction: %v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("the transaction %s does not exist", did)
	}

	var transaction models.Transaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	return &transaction, nil
}

func (s *SmartContract) UpdateTransaction(ctx contractapi.TransactionContextInterface, transactionJSON string) error {
	var transaction models.Transaction
	err := json.Unmarshal([]byte(transactionJSON), &transaction)
	if err != nil {
		return fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	transactionDid := transaction.Did.ToString()
	readTransaction, _ := s.ReadTransaction(ctx, transactionDid)

	if readTransaction == nil {
		return fmt.Errorf("the transaction %s does not exist", transactionDid)
	}

	return ctx.GetStub().PutState(transactionDid, []byte(transactionJSON))
}

func (s *SmartContract) DeleteTransaction(ctx contractapi.TransactionContextInterface, did string) error {
	_, err := s.ReadTransaction(ctx, did)
	if err != nil {
		return fmt.Errorf("the transaction %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}
