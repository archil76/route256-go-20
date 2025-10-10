package postgres

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) GetStock(ctx context.Context, sku int64) (*model.Stock, error) {
	return r.getStock(ctx, sku)
}

func (r *Repository) getStock(ctx context.Context, sku int64) (*model.Stock, error) {
	if sku < 1 {
		return nil, model.ErrSkuIsNotValid
	}
	pool := r.pooler.PickPool(ctx)

	const query = `select id, total_count, reserved from stocks where id = $1`

	rows, err := pool.Query(ctx, query, sku)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rows.Next()
	var stock model.Stock
	if err := rows.Scan(&stock.Sku, &stock.TotalCount, &stock.Reserved); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &stock, nil
}
