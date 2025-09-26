package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) ReserveRemove(ctx context.Context, items []model.Item) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stocks := []model.Stock{}
	for _, item := range items {
		upStock, err := r.getStock(ctx, item.Sku)
		if err != nil {
			return model.ErrStockDoesntExist
		}

		stocks = append(stocks, model.Stock{
			Sku:        upStock.Sku,
			TotalCount: upStock.TotalCount - item.Count,
			Reserved:   upStock.Reserved - item.Count,
		})
	}

	for _, stock := range stocks {
		_, err := r.updateStock(ctx, stock)
		if err != nil {
			return err
		}
	}

	return nil
}
