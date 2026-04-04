package server

import (
	"context"
	commentspb "route256/comments/internal/api"
	"route256/comments/internal/domain/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) Add(ctx context.Context, request *commentspb.CommentAddRequest) (*commentspb.CommentAddResponse, error) {
	comment := model.Comment{}
	comment.UserID = request.UserID
	comment.Comment = request.Comment
	comment.Sku = request.Sku

	commentID, err := s.commentsService.Add(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &commentspb.CommentAddResponse{
		ID: commentID,
	}, status.Errorf(codes.OK, "")
}
