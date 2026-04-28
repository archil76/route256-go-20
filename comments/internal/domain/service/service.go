package comments

import (
	"context"
	"route256/comments/internal/domain/model"
	"time"
)

type Repository interface {
	Add(ctx context.Context, comment model.Comment) (*model.Comment, error)
	Edit(ctx context.Context, order model.Comment) (*model.Comment, error)
	GetByID(ctx context.Context, ID int64) (*model.Comment, error)
	GetListBySKU(ctx context.Context, sku int64) (*model.CommentsList, error)
	GetListByUserID(ctx context.Context, userID int64) (*model.CommentsList, error)
}

type PgPooler interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type CommentsService struct {
	repository   Repository
	pooler       PgPooler
	editInterval time.Duration
}

func NewCommentsService(repository Repository, pooler PgPooler, editInterval time.Duration) *CommentsService {
	return &CommentsService{repository: repository, pooler: pooler, editInterval: editInterval}
}
