package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gzttcydxx/newapi/models"
	"github.com/gzttcydxx/newapi/utils"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// CRUDHandler 通用CRUD处理器
type CRUDHandler[T any] struct {
	ResourceName string
	Contract     *client.Contract
	Methods      models.CRUDMethods
}

// NewCRUDHandler 创建CRUD处理器
func NewCRUDHandler[T any](resourceName string, contract *client.Contract, methods models.CRUDMethods) *CRUDHandler[T] {
	return &CRUDHandler[T]{
		ResourceName: resourceName,
		Contract:     contract,
		Methods:      methods,
	}
}

func (h *CRUDHandler[T]) Create(ctx context.Context, input *models.JSONBody[T]) (*models.JSONBody[models.Status], error) {
	log.Printf("Creating %s with input: %v", h.ResourceName, input)
	bytes, err := json.Marshal(input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to marshal data: %v", err))
	}

	_, err = h.Contract.SubmitTransaction(h.Methods.Create, string(bytes))
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to submit transaction: %v", err))
	}

	return &models.JSONBody[models.Status]{
		Body: models.Status{
			Success: true,
		},
	}, nil
}

func (h *CRUDHandler[T]) Get(ctx context.Context, input *models.GetInput) (*models.JSONBody[T], error) {
	result, err := h.Contract.EvaluateTransaction(h.Methods.Read, input.Did)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to evaluate transaction: %v", err))
	}

	var data T
	err = json.Unmarshal(result, &data)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to unmarshal data: %v", err))
	}

	return &models.JSONBody[T]{Body: data}, nil
}

func (h *CRUDHandler[T]) Query(ctx context.Context, input *models.JSONBody[T]) (*models.JSONBody[models.List[T]], error) {
	query := map[string]interface{}{
		"selector": utils.GetNonZeroFields[T](input.Body),
	}

	queryString, err := json.Marshal(query)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to marshal query: %v", err))
	}

	result, err := h.Contract.EvaluateTransaction(h.Methods.Query, string(queryString))
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to evaluate transaction: %v", err))
	}

	var items []T
	err = json.Unmarshal(result, &items)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to unmarshal query result: %v", err))
	}

	return &models.JSONBody[models.List[T]]{Body: models.List[T]{Items: items}}, nil
}

func (h *CRUDHandler[T]) Update(ctx context.Context, input *models.JSONBody[T]) (*models.JSONBody[models.Status], error) {
	bytes, err := json.Marshal(input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to marshal data: %v", err))
	}

	_, err = h.Contract.SubmitTransaction(h.Methods.Update, string(bytes))
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to submit transaction: %v", err))
	}

	return &models.JSONBody[models.Status]{
		Body: models.Status{
			Success: true,
		},
	}, nil
}

func (h *CRUDHandler[T]) Delete(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	_, err := h.Contract.SubmitTransaction(h.Methods.Delete, input.Did)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to submit transaction: %v", err))
	}

	return &models.JSONBody[models.Status]{
		Body: models.Status{
			Success: true,
		},
	}, nil
}
