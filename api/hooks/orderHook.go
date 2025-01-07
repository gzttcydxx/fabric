package hooks

import (
	"context"

	"github.com/gzttcydxx/newapi/models"
)

type OrderHook struct {
	*CRUDHook[models.Order]
}

func (h *OrderHook) AfterCreate(ctx context.Context, input *models.JSONBody[models.Order]) error {
	return nil
}

func (h *OrderHook) AfterUpdate(ctx context.Context, input *models.JSONBody[models.Order]) error {
	return nil
}

func (h *OrderHook) AfterDelete(ctx context.Context, input *models.GetInput) error {
	return nil
}
