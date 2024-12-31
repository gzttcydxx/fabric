package contract

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) ReadProduct(ctx contractapi.TransactionContextInterface, did string) (models.Product, error) {
	var product models.Product
	result, err := ctx.GetStub().GetState(did)
	if err != nil {
		return product, fmt.Errorf("failed to get product info: %v", err)
	}
	if result == nil {
		return product, fmt.Errorf("the product %s does not exist", did)
	}

	err = json.Unmarshal(result, &product)
	if err != nil {
		return product, fmt.Errorf("failed to unmarshal product info: %v", err)
	}

	return product, nil
}

func (s *SmartContract) ReadProducts(ctx contractapi.TransactionContextInterface) ([]*models.Product, error) {
	// 使用键范围查询
	startKey := "PRODUCT_"
	endKey := "PRODUCT_\ufff0" // \ufff0 是比 _ 大的字符

	iterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer iterator.Close()

	var products []*models.Product
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var product models.Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal product: %v", err)
		}
		products = append(products, &product)
	}

	return products, nil
}

func (s *SmartContract) QueryProducts(ctx contractapi.TransactionContextInterface, query string) ([]*models.Product, error) {
	// 执行富查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	var products []*models.Product
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var product models.Product
		err = json.Unmarshal(queryResult.Value, &product)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal product: %v", err)
		}
		products = append(products, &product)
	}

	return products, nil
}

func (s *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, productJson string) error {
	var product models.Product
	err := json.Unmarshal([]byte(productJson), &product)
	if err != nil {
		return fmt.Errorf("failed to unmarshal product: %v", err)
	}

	did := product.Did

	result, _ := s.ReadProduct(ctx, did.ToString())
	if result.Did.SpecificID != "" {
		return fmt.Errorf("the product %s already exists", did)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(productJson))
	if err != nil {
		return fmt.Errorf("failed to put product: %v", err)
	}

	return nil
}

func (s *SmartContract) UpdateProduct(ctx contractapi.TransactionContextInterface, productJson string) error {
	// 解析更新数据
	var updates models.Product
	if err := json.Unmarshal([]byte(productJson), &updates); err != nil {
		return fmt.Errorf("failed to unmarshal updates: %v", err)
	}

	// 读取现有数据
	existing, err := s.ReadProduct(ctx, updates.Did.ToString())
	if err != nil {
		return fmt.Errorf("failed to get product info: %v", err)
	}
	if existing.Did.SpecificID == "" {
		return fmt.Errorf("the product %s does not exist", updates.Did)
	}

	// 使用反射更新非零值字段
	updatedProduct := existing
	v := reflect.ValueOf(&updates).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		newValue := v.Field(i)

		// 跳过零值字段
		if newValue.IsZero() {
			continue
		}

		// 更新非零值字段
		reflect.ValueOf(&updatedProduct).Elem().FieldByName(field.Name).Set(newValue)
	}

	// 保存更新后的数据
	updatedJson, err := json.Marshal(updatedProduct)
	if err != nil {
		return fmt.Errorf("failed to marshal updated product: %v", err)
	}

	if err := ctx.GetStub().PutState(updates.Did.ToString(), updatedJson); err != nil {
		return fmt.Errorf("failed to put updated product: %v", err)
	}

	return nil
}

func (s *SmartContract) DeleteProduct(ctx contractapi.TransactionContextInterface, did string) error {
	readProduct, err := s.ReadProduct(ctx, did)
	if err != nil {
		return fmt.Errorf("failed to get product info: %v", err)
	}
	if readProduct.Did.SpecificID == "" {
		return fmt.Errorf("the product %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}
