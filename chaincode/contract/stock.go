package contract

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) CreateStock(ctx contractapi.TransactionContextInterface, stock models.Stock) error {
	stockDid := stock.Did.ToString()
	readStock, _ := s.ReadStock(ctx, stockDid)

	if readStock != nil {
		return fmt.Errorf("the stock %s already exists", stockDid)
	}

	stockJSON, err := json.Marshal(stock)
	if err != nil {
		return fmt.Errorf("failed to marshal stock: %v", err)
	}

	return ctx.GetStub().PutState(stockDid, stockJSON)
}

func (s *SmartContract) ReadStock(ctx contractapi.TransactionContextInterface, did string) (*models.Stock, error) {
	stockJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, fmt.Errorf("failed to read stock: %v", err)
	}
	if stockJSON == nil {
		return nil, fmt.Errorf("the stock %s does not exist", did)
	}

	var stock models.Stock
	err = json.Unmarshal(stockJSON, &stock)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stock: %v", err)
	}

	return &stock, nil
}

func (s *SmartContract) UpdateStock(ctx contractapi.TransactionContextInterface, stock models.Stock) error {
	stockDid := stock.Did.ToString()
	readStock, _ := s.ReadStock(ctx, stockDid)

	if readStock == nil {
		return fmt.Errorf("the stock %s does not exist", stockDid)
	}

	stockJSON, err := json.Marshal(stock)
	if err != nil {
		return fmt.Errorf("failed to marshal stock: %v", err)
	}

	return ctx.GetStub().PutState(stockDid, stockJSON)
}

func (s *SmartContract) DeleteStock(ctx contractapi.TransactionContextInterface, did string) error {
	_, err := s.ReadStock(ctx, did)
	if err != nil {
		return fmt.Errorf("the stock %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}
