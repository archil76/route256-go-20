package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) UpdateStock(ctx context.Context, stock model.Stock) (*model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.updateStock(ctx, stock)
}

func (r *Repository) updateStock(ctx context.Context, stock model.Stock) (*model.Stock, error) {
	upStock, err := r.getStock(ctx, stock.Sku)
	if err != nil {
		return nil, err
	}

	upStock.TotalCount = stock.TotalCount
	upStock.Reserved = stock.Reserved

	return upStock, nil
}
