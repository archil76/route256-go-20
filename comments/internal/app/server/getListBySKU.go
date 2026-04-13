package server

import (
	"context"
	commentspb "route256/comments/internal/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s Server) CommentListBySKU(ctx context.Context, request *commentspb.CommentListBySKURequest) (*commentspb.CommentListBySKUResponse, error) {
	commentsList, err := s.commentsService.GetListBySKU(ctx, request.Sku)
	if err != nil {
		return nil, err
	}

	response := commentspb.CommentListBySKUResponse{}

	for _, comment := range commentsList.Comments {

		response.Comments = append(response.Comments, &commentspb.CommentsBySKU{
			ID:        comment.ID,
			UserID:    comment.UserID,
			Comment:   comment.Comment,
			CreatedAt: timestamppb.New(comment.CreatedAt),
		})

	}

	return &response, status.Errorf(codes.OK, "")
}
