package postgres

import (
	"context"
	"route256/comments/internal/domain/model"

	"github.com/cockroachdb/errors"
)

func (r *Repository) GetByID(ctx context.Context, ID int64) (*model.Comment, error) {
	return r.getComment(ctx, ID)
}

func (r *Repository) getComment(ctx context.Context, ID int64) (*model.Comment, error) {

	if ID < 1 {
		return nil, model.ErrCommentIDIsNotValid
	}

	for _, pool := range r.sm.GetShards(ctx) {
		const queryOrders = `SELECT id, user_id, sku, comment, created_at FROM comments where id = $1`

		upComment := model.Comment{}

		if err := pool.QueryRow(ctx, queryOrders, ID).
			Scan(&upComment.ID, &upComment.UserID, &upComment.Sku, &upComment.Comment, &upComment.CreatedAt); err == nil {
			return &upComment, nil
		}
	}

	return nil, errors.Wrap(model.ErrCommentDoesntExist, "pgx.QueryRow.Scan")
}
