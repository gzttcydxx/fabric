package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	didModels "github.com/gzttcydxx/did/models"
	"github.com/gzttcydxx/newapi/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type OrderHandler struct {
	Contract    *client.Contract
	CRUDHandler *CRUDHandler[models.Order]
}

func NewOrderHandler(contract *client.Contract) *OrderHandler {
	return &OrderHandler{
		Contract: contract,
		CRUDHandler: NewCRUDHandler[models.Order]("order", contract, models.CRUDMethods{
			Create: "createOrder",
			Read:   "readOrder",
			Update: "updateOrder",
		}),
	}
}

func (h *OrderHandler) PublicOrder(ctx context.Context, input *models.JSONBody[models.Order]) (*models.JSONBody[models.Status], error) {
	order := input.Body
	order.Status = models.Created
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt
	order.SupplyProduct = []models.Product{}
	did, err := didModels.NewDID(fmt.Sprintf("did:order:%d", order.CreatedAt.Unix()))
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to create did: %v", err))
	}
	order.Did = *did

	return h.CRUDHandler.Create(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

func (h *OrderHandler) SupplierConfirmOrder(ctx context.Context, input *models.OrderProductInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, &models.GetInput{Did: input.Did})
	if err != nil {
		return nil, err
	}

	order := body.Body
	if order.Status > models.SupplierConfirmed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is already selected by demander or approved or finished or cancelled", input.Did))
	}

	order.SupplyProduct = append(order.SupplyProduct, input.Body)
	order.Status = models.SupplierConfirmed
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

func (h *OrderHandler) DemanderSelectOrder(ctx context.Context, input *models.OrderProductInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, &models.GetInput{Did: input.Did})
	if err != nil {
		return nil, err
	}

	order := body.Body
	if order.Status != models.SupplierConfirmed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is not confirmed by supplier", input.Did))
	}

	order.Status = models.DemanderSelected
	order.UpdatedAt = time.Now()
	order.ComfirmProduct = input.Body
	order.SupplierDid = input.Body.OrgDid()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

func (h *OrderHandler) SupplierApproveOrder(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, input)
	if err != nil {
		return nil, err
	}

	order := body.Body
	if order.Status != models.DemanderSelected {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is not selected by demander", input.Did))
	}

	order.Status = models.SupplierApproved
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

func (h *OrderHandler) DemanderApproveOrder(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, input)
	if err != nil {
		return nil, err
	}

	order := body.Body
	if order.Status != models.SupplierApproved {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is not approved by supplier", input.Did))
	}

	order.Status = models.Completed
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

func (h *OrderHandler) SupplierCancelOrder(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, input)
	if err != nil {
		return nil, err
	}

	order := body.Body
	if order.Status == models.Completed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is completed", input.Did))
	}

	order.Status = models.SupplierCanceled
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

func (h *OrderHandler) DemanderCancelOrder(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, input)
	if err != nil {
		return nil, err
	}

	order := body.Body
	if order.Status == models.Completed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is completed", input.Did))
	}

	order.Status = models.DemanderCanceled
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}
