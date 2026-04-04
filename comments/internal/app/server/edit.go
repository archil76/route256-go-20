package server

import (
	"context"
	commentspb "route256/comments/internal/api"
	"route256/comments/internal/domain/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) Edit(ctx context.Context, request *commentspb.CommentEditRequest) (*commentspb.CommentEditResponse, error) {
	comment := model.Comment{}
	comment.UserID = request.UserID
	comment.Comment = request.NewComment

	_, err := s.commentsService.Edit(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &commentspb.CommentEditResponse{}, status.Errorf(codes.OK, "")
}
