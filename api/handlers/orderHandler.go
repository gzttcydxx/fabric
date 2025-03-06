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

// 买方创建订单并选择卖家
func (h *OrderHandler) BuyerCreateOrderWithSeller(ctx context.Context, input *models.JSONBody[models.Order]) (*models.JSONBody[models.OrderResponse], error) {
	order := input.Body
	order.Status = models.Created
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt

	// 确保包含卖家DID
	emptyDID := didModels.DID{}
	if order.SellerDid == emptyDID {
		return nil, huma.Error400BadRequest("seller DID is required")
	}

	// 确保包含买家DID
	if order.BuyerDid == emptyDID {
		return nil, huma.Error400BadRequest("buyer DID is required")
	}

	// 确保包含产品
	emptyProduct := models.Product{}
	if order.Product == emptyProduct {
		return nil, huma.Error400BadRequest("product is required")
	}

	// 确保包含数量
	if order.Num <= 0 {
		return nil, huma.Error400BadRequest("num must be greater than 0")
	}

	did, err := didModels.NewDID(fmt.Sprintf("did:order:%d", order.CreatedAt.Unix()))
	if err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("failed to create did: %v", err))
	}
	order.Did = *did

	// 创建订单
	result, err := h.CRUDHandler.Create(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("failed to create order: %v", err))
	}

	return &models.JSONBody[models.OrderResponse]{
		Body: models.OrderResponse{
			Status:   result.Body,
			OrderDid: order.Did.ToString(),
		},
	}, nil
}

// 卖家确认订单
func (h *OrderHandler) SellerConfirmOrder(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, input)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to get order: %v", err))
	}

	order := body.Body
	if order.Status == models.Completed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is completed", input.Did))
	}
	if order.Status != models.Created {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is not in created status", input.Did))
	}

	order.Status = models.SellerConfirmed
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

// 买家确认订单，确认后交易完成
func (h *OrderHandler) BuyerConfirmOrder(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, input)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to get order: %v", err))
	}

	order := body.Body
	if order.Status == models.Completed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is completed", input.Did))
	}
	if order.Status != models.SellerConfirmed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is not confirmed by seller", input.Did))
	}

	// 直接设置为完成状态
	order.Status = models.Completed
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

// 卖家取消订单
func (h *OrderHandler) SellerCancelOrder(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, input)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to get order: %v", err))
	}

	order := body.Body
	if order.Status == models.Completed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is completed", input.Did))
	}
	if order.Status != models.Created {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is not in created status", input.Did))
	}

	order.Status = models.SellerCanceled
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}

// 买家取消订单
func (h *OrderHandler) BuyerCancelOrder(ctx context.Context, input *models.GetInput) (*models.JSONBody[models.Status], error) {
	body, err := h.CRUDHandler.Get(ctx, input)
	if err != nil {
		return nil, huma.Error400BadRequest(fmt.Sprintf("failed to get order: %v", err))
	}

	order := body.Body
	if order.Status == models.Completed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is completed", input.Did))
	}
	if order.Status != models.SellerConfirmed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("order %s is not confirmed by seller", input.Did))
	}

	order.Status = models.BuyerCanceled
	order.UpdatedAt = time.Now()

	return h.CRUDHandler.Update(ctx, &models.JSONBody[models.Order]{
		Body: order,
	})
}
