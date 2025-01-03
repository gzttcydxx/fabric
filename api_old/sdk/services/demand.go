package services

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type DemandService struct {
	contract *client.Contract
}

func NewDemandService(contract *client.Contract) *DemandService {
	return &DemandService{
		contract: contract,
	}
}

// InitTransaction 初始化交易
func (s *DemandService) InitTransaction(demandItemDid string, demandNum string) (*models.ItemDemand, error) {
	result, err := s.contract.EvaluateTransaction("InitTransaction", demandItemDid, demandNum)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %v", err)
	}

	demand := &models.ItemDemand{}
	if err := json.Unmarshal(result, demand); err != nil {
		return nil, fmt.Errorf("failed to unmarshal demand: %v", err)
	}

	return demand, nil
}

// DealTransaction 处理交易
func (s *DemandService) DealTransaction(selfDid string, stock models.Stock, demandNum string) (*models.Transaction, error) {
	stockJSON, err := json.Marshal(stock)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stock: %v", err)
	}

	result, err := s.contract.EvaluateTransaction("DealTransaction", selfDid, string(stockJSON), demandNum)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %v", err)
	}

	transaction := &models.Transaction{}
	if err := json.Unmarshal(result, transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	return transaction, nil
}

// ConfirmTransaction 确认交易
func (s *DemandService) ConfirmTransaction(accept bool, transaction *models.Transaction) error {
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	acceptStr := "false"
	if accept {
		acceptStr = "true"
	}

	_, err = s.contract.SubmitTransaction("ConfirmTransaction", acceptStr, string(transactionJSON))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	return nil
}
