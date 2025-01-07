package contract

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface, orderJSON string) error {
	var order models.Order
	err := json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order: %v", err)
	}

	err = ctx.GetStub().PutState(order.Did.ToString(), []byte(orderJSON))
	if err != nil {
		return fmt.Errorf("failed to put order: %v", err)
	}

	return nil
}

func (s *SmartContract) ReadOrder(ctx contractapi.TransactionContextInterface, orderDid string) (*models.Order, error) {
	orderJSON, err := ctx.GetStub().GetState(orderDid)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	var order models.Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %v", err)
	}
	return &order, nil
}

func (s *SmartContract) QueryOrders(ctx contractapi.TransactionContextInterface, query string) ([]*models.Order, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	var orders []*models.Order
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}
		var order models.Order
		err = json.Unmarshal(queryResponse.Value, &order)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal order: %v", err)
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (s *SmartContract) UpdateOrder(ctx contractapi.TransactionContextInterface, orderJSON string) error {
	var order models.Order
	err := json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order: %v", err)
	}

	readOrder, err := s.ReadOrder(ctx, order.Did.ToString())
	if err != nil {
		return fmt.Errorf("failed to read order: %v", err)
	}
	if readOrder == nil {
		return fmt.Errorf("order %s not found", order.Did.ToString())
	}

	err = ctx.GetStub().PutState(order.Did.ToString(), []byte(orderJSON))
	if err != nil {
		return fmt.Errorf("failed to put order: %v", err)
	}

	return nil
}

func (s *SmartContract) DeleteOrder(ctx contractapi.TransactionContextInterface, orderDid string) error {
	err := ctx.GetStub().DelState(orderDid)
	if err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}
	return nil
}
