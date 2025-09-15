package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) GetStock(ctx context.Context, sku int64) (*model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.getStock(ctx, sku)
}

func (r *Repository) getStock(_ context.Context, sku int64) (*model.Stock, error) {

	if sku < 1 {
		return nil, ErrSkuIsNotValid
	}

	stock, ok := r.storage[sku]
	if !ok {
		return nil, ErrStockDoesntExist
	}

	return &stock, nil
}
