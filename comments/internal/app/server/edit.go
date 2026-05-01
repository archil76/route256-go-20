package server

import (
	"context"
	commentspb "route256/comments/internal/api"
	"route256/comments/internal/domain/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) CommentEdit(ctx context.Context, request *commentspb.CommentEditRequest) (*commentspb.CommentEditResponse, error) {
	comment := model.Comment{}
	comment.UserID = request.UserID
	comment.ID = request.CommentID
	comment.Comment = request.NewComment

	_, err := s.commentsService.Edit(ctx, comment)
	if err != nil {

		if err == model.ErrUserNotAuthor {
			return nil, status.Errorf(codes.PermissionDenied, "")
		}

		if err == model.ErrEditTimeExpired {
			return nil, status.Errorf(codes.FailedPrecondition, "")
		}

		if err == model.ErrCommentDoesntExist {
			return nil, status.Errorf(codes.NotFound, "")
		}

		return nil, err
	}

	return &commentspb.CommentEditResponse{}, status.Errorf(codes.OK, "")
}
