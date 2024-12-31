package contract

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) ReadPart(ctx contractapi.TransactionContextInterface, did string) (models.Part, error) {
	var part models.Part
	result, err := ctx.GetStub().GetState(did)
	if err != nil {
		return part, fmt.Errorf("failed to get part info: %v", err)
	}
	if result == nil {
		return part, fmt.Errorf("the part %s does not exist", did)
	}

	err = json.Unmarshal(result, &part)
	if err != nil {
		return part, fmt.Errorf("failed to unmarshal part info: %v", err)
	}

	return part, nil
}

func (s *SmartContract) ReadParts(ctx contractapi.TransactionContextInterface) ([]*models.Part, error) {
	// 使用键范围查询
	startKey := "PART_"
	endKey := "PART_\ufff0" // \ufff0 是比 _ 大的字符

	iterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer iterator.Close()

	var parts []*models.Part
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var part models.Part
		err = json.Unmarshal(queryResponse.Value, &part)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal part: %v", err)
		}
		parts = append(parts, &part)
	}

	return parts, nil
}

func (s *SmartContract) QueryParts(ctx contractapi.TransactionContextInterface, query string) ([]*models.Part, error) {
	// 执行富查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	var parts []*models.Part
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var part models.Part
		err = json.Unmarshal(queryResult.Value, &part)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal part: %v", err)
		}
		parts = append(parts, &part)
	}

	return parts, nil
}

func (s *SmartContract) CreatePart(ctx contractapi.TransactionContextInterface, partJson string) error {
	var part models.Part
	err := json.Unmarshal([]byte(partJson), &part)
	if err != nil {
		return fmt.Errorf("failed to unmarshal part: %v", err)
	}

	did := part.Did

	result, _ := s.ReadPart(ctx, did.ToString())
	if result.Did.SpecificID != "" {
		return fmt.Errorf("the part %s already exists", did)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(partJson))
	if err != nil {
		return fmt.Errorf("failed to put part: %v", err)
	}

	return nil
}

func (s *SmartContract) UpdatePart(ctx contractapi.TransactionContextInterface, partJson string) error {
	// 解析更新数据
	var updates models.Part
	if err := json.Unmarshal([]byte(partJson), &updates); err != nil {
		return fmt.Errorf("failed to unmarshal updates: %v", err)
	}

	// 读取现有数据
	existing, err := s.ReadPart(ctx, updates.Did.ToString())
	if err != nil {
		return fmt.Errorf("failed to get part info: %v", err)
	}
	if existing.Did.SpecificID == "" {
		return fmt.Errorf("the part %s does not exist", updates.Did)
	}

	// 使用反射更新非零值字段
	updatedPart := existing
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
		reflect.ValueOf(&updatedPart).Elem().FieldByName(field.Name).Set(newValue)
	}

	// 保存更新后的数据
	updatedJson, err := json.Marshal(updatedPart)
	if err != nil {
		return fmt.Errorf("failed to marshal updated part: %v", err)
	}

	if err := ctx.GetStub().PutState(updates.Did.ToString(), updatedJson); err != nil {
		return fmt.Errorf("failed to put updated part: %v", err)
	}

	return nil
}

func (s *SmartContract) DeletePart(ctx contractapi.TransactionContextInterface, did string) error {
	readPart, err := s.ReadPart(ctx, did)
	if err != nil {
		return fmt.Errorf("failed to get part info: %v", err)
	}
	if readPart.Did.SpecificID == "" {
		return fmt.Errorf("the part %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}
