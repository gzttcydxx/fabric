package hooks

import (
	"context"

	"github.com/gzttcydxx/newapi/models"
)

type CRUDHook[T any] struct{}

func (h *CRUDHook[T]) BeforeCreate(ctx context.Context, input *models.JSONBody[T]) error {
	return nil
}

func (h *CRUDHook[T]) AfterCreate(ctx context.Context, input *models.JSONBody[T]) error {
	return nil
}

func (h *CRUDHook[T]) BeforeGet(ctx context.Context, input *models.GetInput) error {
	return nil
}

func (h *CRUDHook[T]) AfterGet(ctx context.Context, input *models.GetInput) error {
	return nil
}

func (h *CRUDHook[T]) BeforeQuery(ctx context.Context, input *models.JSONBody[T]) error {
	return nil
}

func (h *CRUDHook[T]) AfterQuery(ctx context.Context, input *models.JSONBody[T]) error {
	return nil
}

func (h *CRUDHook[T]) BeforeUpdate(ctx context.Context, input *models.JSONBody[T]) error {
	return nil
}

func (h *CRUDHook[T]) AfterUpdate(ctx context.Context, input *models.JSONBody[T]) error {
	return nil
}

func (h *CRUDHook[T]) BeforeDelete(ctx context.Context, input *models.GetInput) error {
	return nil
}

func (h *CRUDHook[T]) AfterDelete(ctx context.Context, input *models.GetInput) error {
	return nil
}
