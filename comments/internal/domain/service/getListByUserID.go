package comments

import (
	"context"
	"route256/comments/internal/domain/model"
)

func (s *CommentsService) GetListByUserID(ctx context.Context, userID int64) (*model.CommentsList, error) {

	upCommentsList, err := s.repository.GetListByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return upCommentsList, nil
}
