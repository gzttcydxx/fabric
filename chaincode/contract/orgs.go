package contract

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) ReadOrg(ctx contractapi.TransactionContextInterface, did string) (models.Org, error) {
	var org models.Org
	result, err := ctx.GetStub().GetState(did)
	if err != nil {
		return org, fmt.Errorf("failed to get org info: %v", err)
	}
	if result == nil {
		return org, fmt.Errorf("the org %s does not exist", did)
	}

	err = json.Unmarshal(result, &org)
	if err != nil {
		return org, fmt.Errorf("failed to unmarshal org info: %v", err)
	}

	return org, nil
}

func (s *SmartContract) ReadOrgs(ctx contractapi.TransactionContextInterface) ([]*models.Org, error) {
	// 使用键范围查询
	startKey := "ORG_"
	endKey := "ORG_\ufff0" // \ufff0 是比 _ 大的字符

	iterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer iterator.Close()

	var orgs []*models.Org
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var org models.Org
		err = json.Unmarshal(queryResponse.Value, &org)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal org: %v", err)
		}
		orgs = append(orgs, &org)
	}

	return orgs, nil
}

func (s *SmartContract) QueryOrgs(ctx contractapi.TransactionContextInterface, query string) ([]*models.Org, error) {
	// 执行富查询
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	var orgs []*models.Org
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var org models.Org
		err = json.Unmarshal(queryResult.Value, &org)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal org: %v", err)
		}
		orgs = append(orgs, &org)
	}

	return orgs, nil
}

func (s *SmartContract) CreateOrg(ctx contractapi.TransactionContextInterface, orgJson string) error {
	var org models.Org
	err := json.Unmarshal([]byte(orgJson), &org)
	if err != nil {
		return fmt.Errorf("failed to unmarshal org: %v", err)
	}

	did := org.Did

	result, _ := s.ReadOrg(ctx, did.ToString())
	if result.Did.SpecificID != "" {
		return fmt.Errorf("the org %s already exists", did)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(orgJson))
	if err != nil {
		return fmt.Errorf("failed to put org: %v", err)
	}

	return nil
}

func (s *SmartContract) UpdateOrg(ctx contractapi.TransactionContextInterface, orgJson string) error {
	// 解析更新数据
	var updates models.Org
	if err := json.Unmarshal([]byte(orgJson), &updates); err != nil {
		return fmt.Errorf("failed to unmarshal updates: %v", err)
	}

	// 读取现有数据
	existing, err := s.ReadOrg(ctx, updates.Did.ToString())
	if err != nil {
		return fmt.Errorf("failed to get org info: %v", err)
	}
	if existing.Did.SpecificID == "" {
		return fmt.Errorf("the org %s does not exist", updates.Did)
	}

	// 使用反射更新非零值字段
	updatedOrg := existing
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
		reflect.ValueOf(&updatedOrg).Elem().FieldByName(field.Name).Set(newValue)
	}

	// 保存更新后的数据
	updatedJson, err := json.Marshal(updatedOrg)
	if err != nil {
		return fmt.Errorf("failed to marshal updated org: %v", err)
	}

	if err := ctx.GetStub().PutState(updates.Did.ToString(), updatedJson); err != nil {
		return fmt.Errorf("failed to put updated org: %v", err)
	}

	return nil
}

func (s *SmartContract) DeleteOrg(ctx contractapi.TransactionContextInterface, did string) error {
	readOrg, err := s.ReadOrg(ctx, did)
	if err != nil {
		return fmt.Errorf("failed to get org info: %v", err)
	}
	if readOrg.Did.SpecificID == "" {
		return fmt.Errorf("the org %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}
