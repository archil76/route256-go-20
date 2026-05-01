package server

import (
	"context"
	commentspb "route256/comments/internal/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s Server) CommentListByUser(ctx context.Context, request *commentspb.CommentListByUserRequest) (*commentspb.CommentListByUserResponse, error) {
	commentsList, err := s.commentsService.GetListByUserID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	response := commentspb.CommentListByUserResponse{}

	for _, comment := range commentsList.Comments {

		response.Comments = append(response.Comments, &commentspb.CommentsByUserID{
			ID:        comment.ID,
			Sku:       comment.Sku,
			Comment:   comment.Comment,
			CreatedAt: timestamppb.New(comment.CreatedAt),
		})

	}

	return &response, status.Errorf(codes.OK, "")
}
