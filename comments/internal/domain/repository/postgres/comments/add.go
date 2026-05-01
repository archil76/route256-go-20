package postgres

import (
	"context"
	"route256/comments/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) Add(ctx context.Context, comment model.Comment) (*model.Comment, error) {
	var err error
	shardIndex := r.sm.GetShardIndex(comment.Sku)
	pool, err := r.sm.PickPool(ctx, shardIndex)
	if err != nil {
		return nil, err
	}

	const queryOrders = `INSERT INTO comments (id, user_id, sku, comment) VALUES (nextval('comments_id_manual_seq') + $1, $2, $3, $4) returning id`
	if err = pool.QueryRow(ctx, queryOrders, shardIndex, comment.UserID, comment.Sku, comment.Comment).Scan(&comment.ID); err != nil {
		err = errors.Wrap(err, "pgx.QueryRow.Scan")
		return nil, err
	}

	return &comment, nil
}
