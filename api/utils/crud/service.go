package crud

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/api/utils/reflect"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// CRUDMethods 定义 CRUD 方法名
type CRUDMethods struct {
	Create string // 创建方法名
	Read   string // 读取方法名
	Query  string // 查询方法名
	Update string // 更新方法名
	Delete string // 删除方法名
}

// CRUDService 通用 CRUD 服务
type CRUDService[T any] struct {
	Contract *client.Contract
	Type     string      // 如 "org", "part", "product" 等，用于 did method
	Methods  CRUDMethods // 方法名配置
}

// NewCRUDService 创建新的 CRUD 服务
func NewCRUDService[T any](contract *client.Contract, typeName string, methods CRUDMethods) *CRUDService[T] {
	return &CRUDService[T]{
		Contract: contract,
		Type:     typeName,
		Methods:  methods,
	}
}

// Create 创建记录
func (s *CRUDService[T]) Create(data T) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	_, err = s.Contract.SubmitTransaction(s.Methods.Create, string(bytes))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	return nil
}

// Read 读取记录
func (s *CRUDService[T]) Read(did string) (T, error) {
	var data T
	result, err := s.Contract.EvaluateTransaction(s.Methods.Read, did)
	if err != nil {
		return data, fmt.Errorf("failed to evaluate transaction: %v", err)
	}
	err = json.Unmarshal(result, &data)
	if err != nil {
		return data, fmt.Errorf("failed to unmarshal did doc: %v", err)
	}
	return data, nil
}

// Query 查询记录
func (s *CRUDService[T]) Query(params T) ([]T, error) {
	query := map[string]interface{}{
		"selector": reflect.GetNonZeroFields(params),
	}

	queryString, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %v", err)
	}

	result, err := s.Contract.EvaluateTransaction(s.Methods.Query, string(queryString))
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %v", err)
	}

	var items []T
	err = json.Unmarshal(result, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal query result: %v", err)
	}
	return items, nil
}

// Update 更新记录
func (s *CRUDService[T]) Update(data T) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	_, err = s.Contract.SubmitTransaction(s.Methods.Update, string(bytes))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	return nil
}

// Delete 删除记录
func (s *CRUDService[T]) Delete(did string) error {
	_, err := s.Contract.SubmitTransaction(s.Methods.Delete, did)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}
	return nil
}
