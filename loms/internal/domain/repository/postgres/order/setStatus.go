package inmemoryrepository

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) SetStatus(ctx context.Context, order model.Order, status model.Status) error {

	const query = `UPDATE orders SET status=$2 where id = $1`
	if _, err := r.pool.Exec(ctx, query, order.OrderID, status); err != nil {
		return errors.Wrap(err, "pgx.QueryRow.Scan")
	}

	return nil
}
