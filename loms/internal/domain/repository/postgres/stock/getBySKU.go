package postgres

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) GetBySKU(ctx context.Context, sku int64) (uint32, error) {
	const query = `select total_count-reserved from stocks where id = $1`

	rows, err := r.pool.Query(ctx, query, sku)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	rows.Next()
	var count uint32
	if err := rows.Scan(&count); err != nil {
		return 0, model.ErrStockDoesntExist
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
