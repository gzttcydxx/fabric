package contract

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/gzttcydxx/fabric/chaincode/models"
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

func (s *SmartContract) DeleteAllTransactions(ctx contractapi.TransactionContextInterface) error {
	// 获取所有以 did:transaction: 开头的键
	queryString := `{
        "selector": {
            "_id": {
                "$regex": "^did:transaction:"
            }
        }
    }`

	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return fmt.Errorf("failed to get transactions: %v", err)
	}
	defer iterator.Close()

	var keys []string
	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return fmt.Errorf("failed to iterate transactions: %v", err)
		}
		keys = append(keys, kv.Key)
	}

	// 排序确保顺序一致
	sort.Strings(keys)

	// 批量删除
	for _, key := range keys {
		err = ctx.GetStub().DelState(key)
		if err != nil {
			return fmt.Errorf("failed to delete transaction %s: %v", key, err)
		}
	}

	log.Printf("Successfully deleted %d transactions", len(keys))

	return nil
}
