package postgres

import (
	"context"
	"route256/comments/internal/domain/model"
)

func (r *Repository) Edit(ctx context.Context, comment model.Comment) (*model.Comment, error) {
	return r.edit(ctx, comment)
}

func (r *Repository) edit(ctx context.Context, comment model.Comment) (*model.Comment, error) {
	const query = `UPDATE comments SET comment=$2 where id = $1 returning id`

	upComment := model.Comment{}

	for _, pool := range r.sm.GetShards(ctx) {
		if err := pool.QueryRow(ctx, query, comment.ID, comment.Comment).
			Scan(&upComment.ID); err == nil {
			upComment.Comment = comment.Comment
			return &upComment, nil
		}
	}

	return nil, model.ErrCommentDoesntExist
}
