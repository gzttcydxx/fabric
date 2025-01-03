package services

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type SupplyService struct {
	contract *client.Contract
}

func NewSupplyService(contract *client.Contract) *SupplyService {
	return &SupplyService{
		contract: contract,
	}
}

// GetAvailableStock 获取可用库存
func (s *SupplyService) GetAvailableStock(itemDemand *models.ItemDemand) (*models.Stock, error) {
	itemDemandJSON, err := json.Marshal(itemDemand)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal item demand: %v", err)
	}

	result, err := s.contract.EvaluateTransaction("GetAvailableStock", string(itemDemandJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to get available stock: %v", err)
	}

	var stock *models.Stock
	if err := json.Unmarshal(result, &stock); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stock: %v", err)
	}

	return stock, nil
}

// AcceptTransaction 接受或拒绝交易
func (s *SupplyService) AcceptTransaction(accept bool, transaction *models.Transaction) error {
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	acceptStr := "false"
	if accept {
		acceptStr = "true"
	}

	_, err = s.contract.SubmitTransaction("AcceptTransaction", acceptStr, string(transactionJSON))
	if err != nil {
		return fmt.Errorf("failed to accept/reject transaction: %v", err)
	}

	return nil
}
