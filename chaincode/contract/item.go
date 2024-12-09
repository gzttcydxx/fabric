package contract

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/chaincode/models"
	didModels "github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//////////////////// ItemType 的增删改查 - START ////////////////////

func (s *SmartContract) CreateItemType(ctx contractapi.TransactionContextInterface, itemTypeJSON string) error {
	var itemType models.ItemType
	err := json.Unmarshal([]byte(itemTypeJSON), &itemType)
	if err != nil {
		return fmt.Errorf("failed to unmarshal itemType: %v", err)
	}

	itemTypeDid := itemType.Did.ToString()
	readItemType, _ := s.ReadItemType(ctx, itemTypeDid)

	if readItemType != nil {
		return fmt.Errorf("the itemType %s already exists", itemTypeDid)
	}

	return ctx.GetStub().PutState(itemTypeDid, []byte(itemTypeJSON))
}

func (s *SmartContract) ReadItemType(ctx contractapi.TransactionContextInterface, did string) (*models.ItemType, error) {
	itemTypeJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, fmt.Errorf("failed to read itemType: %v", err)
	}
	if itemTypeJSON == nil {
		return nil, fmt.Errorf("the itemType %s does not exist", did)
	}

	var itemType models.ItemType
	err = json.Unmarshal(itemTypeJSON, &itemType)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal itemType: %v", err)
	}

	return &itemType, nil
}

func (s *SmartContract) UpdateItemType(ctx contractapi.TransactionContextInterface, itemTypeJSON string) error {
	var itemType models.ItemType
	err := json.Unmarshal([]byte(itemTypeJSON), &itemType)
	if err != nil {
		return fmt.Errorf("failed to unmarshal itemType: %v", err)
	}

	itemTypeDid := itemType.Did.ToString()
	readItemType, _ := s.ReadItemType(ctx, itemTypeDid)

	if readItemType == nil {
		return fmt.Errorf("the itemType %s does not exist", itemTypeDid)
	}

	return ctx.GetStub().PutState(itemTypeDid, []byte(itemTypeJSON))
}

func (s *SmartContract) DeleteItemType(ctx contractapi.TransactionContextInterface, did string) error {
	_, err := s.ReadItemType(ctx, did)
	if err != nil {
		return fmt.Errorf("the itemType %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}

//////////////////// ItemType 的增删改查 - END ////////////////////

//////////////////// Item 的增删改查 - START ////////////////////

func (s *SmartContract) CreateItem(ctx contractapi.TransactionContextInterface, itemJSON string) error {
	var item models.Item
	err := json.Unmarshal([]byte(itemJSON), &item)
	if err != nil {
		return fmt.Errorf("failed to unmarshal item: %v", err)
	}

	itemDid := item.Did.ToString()
	readItem, _ := s.ReadItem(ctx, itemDid)

	if readItem != nil {
		return fmt.Errorf("the item %s already exists", itemDid)
	}

	return ctx.GetStub().PutState(itemDid, []byte(itemJSON))
}

func (s *SmartContract) ReadItem(ctx contractapi.TransactionContextInterface, did string) (*models.Item, error) {
	itemJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, fmt.Errorf("failed to read item: %v", err)
	}
	if itemJSON == nil {
		return nil, fmt.Errorf("the item %s does not exist", did)
	}

	var item models.Item
	err = json.Unmarshal(itemJSON, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %v", err)
	}

	return &item, nil
}

func (s *SmartContract) UpdateItem(ctx contractapi.TransactionContextInterface, itemJSON string) error {
	var item models.Item
	err := json.Unmarshal([]byte(itemJSON), &item)
	if err != nil {
		return fmt.Errorf("failed to unmarshal item: %v", err)
	}

	itemDid := item.Did.ToString()
	readItem, _ := s.ReadItem(ctx, itemDid)

	if readItem == nil {
		return fmt.Errorf("the item %s does not exist", itemDid)
	}

	return ctx.GetStub().PutState(itemDid, []byte(itemJSON))
}

func (s *SmartContract) DeleteItem(ctx contractapi.TransactionContextInterface, did string) error {
	_, err := s.ReadItem(ctx, did)
	if err != nil {
		return fmt.Errorf("the item %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}

//////////////////// Item 的增删改查 - END ////////////////////

// 从 ItemType 中找到所有符合条件的 Item 的库存
func (s *SmartContract) FindItemStock(ctx contractapi.TransactionContextInterface, itemTypeDid string) ([]*models.ItemStock, error) {
	var typeDID didModels.DID
	typeDID.FromString(itemTypeDid)

	queryString := fmt.Sprintf(`{
        "selector": {
            "type.did": {
                "Scheme": "%s",
                "Method": "%s",
                "SpecificID": "%s"
            }
        }
    }`, typeDID.Scheme, typeDID.Method, typeDID.SpecificID)

	var items []*models.ItemStock
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get query items: %v", err)
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		itemKey, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next item: %v", err)
		}

		var item models.ItemStock
		err = json.Unmarshal(itemKey.Value, &item)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %v", err)
		}

		items = append(items, &item)
	}

	return items, nil
}
