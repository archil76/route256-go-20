package inmemoryrepository

import (
	"context"
)

func (r *Repository) GetBySKU(ctx context.Context, sku int64) (uint32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stock, err := r.getStock(ctx, sku)
	if err != nil {
		return 0, err
	}

	return stock.TotalCount - stock.Reserved, nil
}
