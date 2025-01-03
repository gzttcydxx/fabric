package contract

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	didModels "github.com/gzttcydxx/did/models"
	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// InitTransaction 初始化交易。
// Parameters:
//   - ctx: 交易上下文接口，用于访问链码的状态和交易信息
//   - demandItemDid: 需求的对象类型DID
//   - demandNum：需求的数量
//
// Returns:
//   - *models.ItemDemand: 返回需求对象，包含需求的对象类型名称、单位等信息
//   - error: 错误信息
func (s *SmartContract) InitTransaction(ctx contractapi.TransactionContextInterface, demandItemDid string, demandNum string) (*models.ItemDemand, error) {
	// 参数验证
	if len(demandItemDid) == 0 {
		return nil, fmt.Errorf("demandItemDid cannot be empty")
	}
	if len(demandNum) == 0 {
		return nil, fmt.Errorf("demandNum cannot be empty")
	}
	fmt.Println("1")

	// 利用 item.Did 从区块链中读取 ItemType 对象
	demandItemJSON, err := ctx.GetStub().GetState(demandItemDid)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if demandItemJSON == nil {
		return nil, fmt.Errorf("demand item type %s does not exist", demandItemDid)
	}
	fmt.Println("2")
	var itemType models.ItemType
	err = json.Unmarshal(demandItemJSON, &itemType)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item type: %v", err)
	}
	fmt.Println("3")
	// 转换数量为整数
	num, err := strconv.Atoi(demandNum)
	if err != nil {
		return nil, fmt.Errorf("invalid num: %v", err)
	}

	// 创建新的 ItemDemand 对象
	demand := &models.ItemDemand{
		ItemType: itemType, // 使用指针
		Num:      num,
	}
	fmt.Println(demand.ItemType)
	return demand, nil
}

// DealTransaction 处理交易请求，根据指定的策略从库存中选择合适的零件。
// Parameters:
//   - ctx: 交易上下文接口，用于访问链码的状态和交易信息
//   - stockJSON: 库存的JSON字符串
//   - demandNum: 需求数量的字符串
//   - method: 可选参数，指定选择策略，默认为"min"
//
// Returns:
//   - *models.Item: 选择的零件
//   - error: 错误信息
func (s *SmartContract) DealTransaction(ctx contractapi.TransactionContextInterface, selfDid string, stockJSON string, demandNum string) (*models.Transaction, error) {
	strategy := "min"

	var stock models.Stock
	err := json.Unmarshal([]byte(stockJSON), &stock)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stock: %v", err)
	}

	num, err := strconv.Atoi(demandNum)
	if err != nil {
		return nil, fmt.Errorf("invalid num: %v", err)
	}

	if strategy == "min" {
		// 从多个匹配的库存中，选择其中价格最低的零件。
		// 如果库存中没有匹配的零件，返回错误。
		var selectedItemStock *models.ItemStock
		lowestPrice := int(^uint(0) >> 1) // 最大整数值

		for _, itemStock := range stock.Items {
			price := itemStock.Item.Price
			if price < lowestPrice && itemStock.Num >= num {
				lowestPrice = price
				selectedItemStock = itemStock
			}
		}

		if selectedItemStock == nil {
			return nil, fmt.Errorf("no enough item found")
		}

		selfDID := didModels.DID{}
		selfDID.FromString(selfDid)
		transactionDid := didModels.DID{}
		transactionDid.FromString("did:example:transaction")

		transaction := &models.Transaction{
			Did:    transactionDid,
			Demand: selfDID,
			Supply: selectedItemStock.Item.Did,
			Type:   "buy",
			Amount: num,
			Item:   selectedItemStock.Item,
			Time:   time.Now(),
			Status: models.START,
		}

		return transaction, nil

	} else {
		return nil, fmt.Errorf("invalid strategy: %s", strategy)
	}
}

func (s *SmartContract) ConfirmTransaction(ctx contractapi.TransactionContextInterface, accept string, transactionJSON string) error {
	var transaction models.Transaction
	err := json.Unmarshal([]byte(transactionJSON), &transaction)
	if err != nil {
		return fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	transactionDid := transaction.Did.ToString()
	if accept != "true" {
		transaction.Status = models.REJECT
	} else {
		transaction.Status = models.END
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
