package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) Reserve(ctx context.Context, items []model.Item) ([]model.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var stocks []model.Stock
	for _, item := range items {
		upStock, err := r.getStock(ctx, item.Sku)
		if err != nil {
			return nil, ErrStockDoesntExist
		}

		available := upStock.TotalCount - upStock.Reserved
		if available < item.Count {
			return nil, ErrShortOfStock
		}

		stocks = append(stocks, model.Stock{
			Sku:        upStock.Sku,
			TotalCount: upStock.TotalCount,
			Reserved:   upStock.Reserved + item.Count,
		})
	}

	for _, stock := range stocks {
		_, err := r.updateStock(ctx, stock)
		if err != nil {
			return nil, err
		}
	}

	return stocks, nil
}
