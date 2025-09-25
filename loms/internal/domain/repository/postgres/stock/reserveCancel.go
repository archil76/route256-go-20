package postgres

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) ReserveCancel(ctx context.Context, items []model.Item) error {
	return nil
}
