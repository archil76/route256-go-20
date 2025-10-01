package postgres

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) UpdateStock(ctx context.Context, stock model.Stock) (*model.Stock, error) {
	return r.updateStock(ctx, stock)
}

func (r *Repository) updateStock(ctx context.Context, stock model.Stock) (*model.Stock, error) {
	const query = `UPDATE stocks SET total_count=$2, reserved=$3 where id = $1`

	upStock := model.Stock{}
	if err := r.pool.QueryRow(ctx, query, stock.Sku, stock.TotalCount, stock.Reserved).
		Scan(&upStock.Sku); err != nil {
		return nil, errors.Wrap(err, "pgx.QueryRow.Scan")
	}
	upStock.TotalCount = stock.TotalCount
	upStock.Reserved = stock.Reserved

	return &upStock, nil
}
