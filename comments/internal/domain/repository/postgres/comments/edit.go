package postgres

import (
	"context"
	"route256/comments/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) Edit(ctx context.Context, comment model.Comment) (*model.Comment, error) {
	return r.edit(ctx, comment)
}

func (r *Repository) edit(ctx context.Context, comment model.Comment) (*model.Comment, error) {
	pool, err := r.sm.PickPool(ctx, comment.Sku)
	if err != nil {
		return nil, err
	}

	const query = `UPDATE comments SET comment=$2 where id = $1 returning id`

	upComment := model.Comment{}
	if err := pool.QueryRow(ctx, query, comment.ID, comment.Comment).
		Scan(&upComment.ID); err != nil {
		return nil, errors.Wrap(err, "pgx.QueryRow.Scan")
	}
	upComment.Comment = comment.Comment

	return &upComment, nil
}
