package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) SetStatus(ctx context.Context, order model.Order, status model.Status) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order.Status = status
	_, err := r.UpdateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
