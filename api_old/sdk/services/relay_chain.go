package services

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type RelayChainService struct {
	contract *client.Contract
}

func NewRelayChainService(contract *client.Contract) *RelayChainService {
	return &RelayChainService{
		contract: contract,
	}
}

// FilterDemandStock 根据需求筛选库存
func (s *RelayChainService) FilterDemandStock(itemDemand *models.ItemDemand) (*models.Stock, error) {
	itemDemandJSON, err := json.Marshal(itemDemand)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal item demand: %v", err)
	}

	result, err := s.contract.EvaluateTransaction("FilterDemandStock", string(itemDemandJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to filter demand stock: %v", err)
	}

	var stock *models.Stock
	if err := json.Unmarshal(result, &stock); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stock: %v", err)
	}

	return stock, nil
}

// Send2SupplyTransaction 发送交易到供给链
func (s *RelayChainService) Send2SupplyTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	result, err := s.contract.SubmitTransaction("Send2SupplyTransaction", string(transactionJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction to supply chain: %v", err)
	}

	var updatedTransaction *models.Transaction
	if err := json.Unmarshal(result, &updatedTransaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	return updatedTransaction, nil
}

// Send2DemandTransaction 发送交易到需求链
func (s *RelayChainService) Send2DemandTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	result, err := s.contract.SubmitTransaction("Send2DemandTransaction", string(transactionJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction to demand chain: %v", err)
	}

	var updatedTransaction *models.Transaction
	if err := json.Unmarshal(result, &updatedTransaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	return updatedTransaction, nil
}
