package postgres

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) Reserve(ctx context.Context, items []model.Item) ([]model.Stock, error) {
	return nil, nil
}
