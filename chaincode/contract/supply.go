package contract

import (
	"encoding/json"
	"fmt"

	didModels "github.com/gzttcydxx/did/models"
	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) GetAvailableStock(ctx contractapi.TransactionContextInterface, itemDemandJSON string) (*models.Stock, error) {
	var itemDemand models.ItemDemand
	err := json.Unmarshal([]byte(itemDemandJSON), &itemDemand)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item demand: %v", err)
	}

	var stock *models.Stock
	items, err := s.FindItemStock(ctx, itemDemand.Did.ToString())
	if err != nil {
		return nil, fmt.Errorf("failed to find items: %v", err)
	}

	var stockDID didModels.DID
	stockDID.FromString("did:example:stock_supply")
	stock = &models.Stock{
		Did:   stockDID,
		Items: make(map[string]*models.ItemStock),
	}

	for _, item := range items {
		if item.Num >= itemDemand.Num {
			stock.Items[item.Did.ToString()] = item
		}
	}

	return stock, nil
}

func (s *SmartContract) AcceptTransaction(ctx contractapi.TransactionContextInterface, accept string, transactionJSON string) error {
	var transaction models.Transaction
	err := json.Unmarshal([]byte(transactionJSON), &transaction)
	if err != nil {
		return fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	transactionDid := transaction.Did.ToString()
	if accept != "true" {
		transaction.Status = models.REJECT
	} else {
		transaction.Status = models.ACCEPT
	}

	transactionJSONbyte, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// 记录交易信息
	err = ctx.GetStub().PutState(transactionDid, transactionJSONbyte)
	if err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}
	return nil
}
