package server

import (
	"context"
	desc "route256/comments/internal/api"
	"route256/comments/internal/domain/model"
)

type CommentsService interface {
	Add(ctx context.Context, comment model.Comment) (int64, error)
	Edit(ctx context.Context, order model.Comment) (int64, error)
	GetByID(ctx context.Context, ID int64) (*model.Comment, error)
	GetListBySKU(ctx context.Context, sku int64) (*model.CommentsList, error)
	GetListByUserID(ctx context.Context, userID int64) (*model.CommentsList, error)
}

type Server struct {
	commentsService CommentsService
	desc.UnimplementedCommentsServer
}

func NewServer(commentsService CommentsService) *Server {
	return &Server{commentsService: commentsService}
}
