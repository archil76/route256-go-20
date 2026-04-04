package postgres

import (
	"context"
	"route256/comments/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) Add(ctx context.Context, comment model.Comment) (*model.Comment, error) {
	var err error

	pool, err := r.sm.PickPool(ctx, comment.Sku)
	if err != nil {
		return nil, err
	}

	const queryOrders = `INSERT INTO comments (user_id, sku, comment) VALUES ($1, $2, $3) returning id`
	if err = pool.QueryRow(ctx, queryOrders, comment.UserID, comment.Sku, comment.Comment).Scan(&comment.ID); err != nil {
		err = errors.Wrap(err, "pgx.QueryRow.Scan")
		return nil, err
	}

	return &comment, nil
}
