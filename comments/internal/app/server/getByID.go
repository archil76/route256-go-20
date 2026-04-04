package server

import (
	"context"
	commentspb "route256/comments/internal/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s Server) GetById(ctx context.Context, request *commentspb.CommentGetByIDRequest) (*commentspb.CommentGetByIDResponse, error) {
	comment, err := s.commentsService.GetByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &commentspb.CommentGetByIDResponse{
		UserID:    comment.UserID,
		Sku:       comment.Sku,
		Comment:   comment.Comment,
		ID:        comment.ID,
		CreatedAt: timestamppb.New(comment.CreatedAt),
	}, status.Errorf(codes.OK, "")
}
