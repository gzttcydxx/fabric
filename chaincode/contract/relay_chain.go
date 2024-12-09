package contract

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/chaincode/models"
	didModels "github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// FilterDemandStock 根据需求筛选库存中的物品
// Parameters:
//   - ctx: 交易上下文接口
//   - demandDid: 需求的DID
//   - itemDemandJSON: 需求物品的JSON字符串
//
// Returns:
//   - *models.Stock: 符合需求的库存
//   - error: 错误信息
func (s *SmartContract) FilterDemandStock(ctx contractapi.TransactionContextInterface, itemDemandJSON string) (*models.Stock, error) {
	// 库存样例
	stock := models.Stock{
		Did: didModels.DID{
			Scheme:     "did",
			Method:     "stock",
			ChainID:    "",
			SpecificID: "1",
			Fragment:   "",
		},
		Items: map[string]*models.ItemStock{
			"did:item:1": {
				Item: models.Item{
					Did: didModels.DID{
						Scheme:     "did",
						Method:     "item",
						ChainID:    "",
						SpecificID: "1",
						Fragment:   "",
					},
					Name: "iPhone 13",
					Type: models.ItemType{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "type",
							ChainID:    "",
							SpecificID: "1",
							Fragment:   "",
						},
						Name: "Smartphone",
						Unit: "piece",
					},
					Org: models.Org{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "org",
							ChainID:    "",
							SpecificID: "1",
							Fragment:   "",
						},
						Name: "Apple Store",
					},
					Owner: models.User{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "user",
							ChainID:    "",
							SpecificID: "1",
							Fragment:   "",
						},
						Name: "John Doe",
						Role: "store_manager",
						Org: models.Org{
							Did: didModels.DID{
								Scheme:     "did",
								Method:     "org",
								ChainID:    "",
								SpecificID: "1",
								Fragment:   "",
							},
							Name: "Apple Store",
						},
					},
					Price: 999,
				},
				Num: 50,
			},
			"did:item:2": {
				Item: models.Item{
					Did: didModels.DID{
						Scheme:     "did",
						Method:     "item",
						ChainID:    "",
						SpecificID: "2",
						Fragment:   "",
					},
					Name: "MacBook Pro",
					Type: models.ItemType{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "type",
							ChainID:    "",
							SpecificID: "2",
							Fragment:   "",
						},
						Name: "Laptop",
						Unit: "piece",
					},
					Org: models.Org{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "org",
							ChainID:    "",
							SpecificID: "1",
							Fragment:   "",
						},
						Name: "Apple Store",
					},
					Owner: models.User{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "user",
							ChainID:    "",
							SpecificID: "1",
							Fragment:   "",
						},
						Name: "John Doe",
						Role: "store_manager",
						Org: models.Org{
							Did: didModels.DID{
								Scheme:     "did",
								Method:     "org",
								ChainID:    "",
								SpecificID: "1",
								Fragment:   "",
							},
							Name: "Apple Store",
						},
					},
					Price: 1999,
				},
				Num: 30,
			},
			"did:item:3": {
				Item: models.Item{
					Did: didModels.DID{
						Scheme:     "did",
						Method:     "item",
						ChainID:    "",
						SpecificID: "3",
						Fragment:   "",
					},
					Name: "AirPods Pro",
					Type: models.ItemType{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "type",
							ChainID:    "",
							SpecificID: "3",
							Fragment:   "",
						},
						Name: "Earphones",
						Unit: "piece",
					},
					Org: models.Org{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "org",
							ChainID:    "",
							SpecificID: "1",
							Fragment:   "",
						},
						Name: "Apple Store",
					},
					Owner: models.User{
						Did: didModels.DID{
							Scheme:     "did",
							Method:     "user",
							ChainID:    "",
							SpecificID: "1",
							Fragment:   "",
						},
						Name: "John Doe",
						Role: "store_manager",
						Org: models.Org{
							Did: didModels.DID{
								Scheme:     "did",
								Method:     "org",
								ChainID:    "",
								SpecificID: "1",
								Fragment:   "",
							},
							Name: "Apple Store",
						},
					},
					Price: 249,
				},
				Num: 100,
			},
		},
	}

	// 解析需求类型
	var itemDemand models.ItemDemand
	err := json.Unmarshal([]byte(itemDemandJSON), &itemDemand)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal itemDemand: %v", err)
	}
	demandDid := itemDemand.Did
	demandType := itemDemand.ItemType

	// 身份认证（忽略细节）
	demandDidJSON, err := ctx.GetStub().GetState(demandDid.ToString())
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if demandDidJSON == nil {
		return nil, fmt.Errorf("identity %s does not exist", demandDid)
	}

	// 筛选符合条件的类型
	demandStock := make(map[string]*models.ItemStock)
	for _, itemStock := range stock.Items {
		if itemStock.Item.Type.Did == demandType.Did {
			demandStock[itemStock.Item.Did.ToString()] = itemStock
		}
	}

	// 返回库存
	return &models.Stock{
		Did: didModels.DID{
			Scheme:     "did",
			Method:     "stock",
			SpecificID: "return",
		},
		Items: demandStock,
	}, nil
}

// TODO:
// 当需求链发起需求后，将需求存入中继链中，并向中继链客户端发起事件
// 中继链客户端向所有拥有该类型item的区块链发送RESTful请求
// 所有满足需求数量的区块链均返回可供给库存，并由中继链返回给需求链

// Send2SupplyTransaction 发送交易到供给链
// Parameters:
//   - ctx: 交易上下文接口
//   - transactionJSON: 交易信息的JSON字符串
//
// Returns:
//   - *models.Transaction: 生成的交易对象指针
//   - error: 错误信息
func (s *SmartContract) Send2SupplyTransaction(ctx contractapi.TransactionContextInterface, transactionJSON string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := json.Unmarshal([]byte(transactionJSON), &transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	transactionDid := transaction.Did.ToString()
	transaction.Status = models.PADING
	transactionJSONbyte, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// 记录交易信息
	err = ctx.GetStub().PutState(transactionDid, transactionJSONbyte)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	// 生成交易
	return &transaction, nil
}

// Send2DemandTransaction 发送交易到需求链
// Parameters:
//   - ctx: 交易上下文接口
//   - transactionJSON: 交易信息的JSON字符串
//
// Returns:
//   - *models.Transaction: 交易对象指针
//   - error: 错误信息
func (s *SmartContract) Send2DemandTransaction(ctx contractapi.TransactionContextInterface, transactionJSON string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := json.Unmarshal([]byte(transactionJSON), &transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	transactionDid := transaction.Did.ToString()
	if transaction.Status != models.REJECT {
		transaction.Status = models.END
	}
	transactionJSONbyte, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// 更新交易信息
	err = ctx.GetStub().PutState(transactionDid, transactionJSONbyte)
	if err != nil {
		return nil, fmt.Errorf("failed to put state: %v", err)
	}

	// 返回交易
	return &transaction, nil
}
