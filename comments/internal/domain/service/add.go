package comments

import (
	"context"
	"route256/comments/internal/domain/model"
)

func (s *CommentsService) Add(ctx context.Context, comment model.Comment) (int64, error) {

	upComment, err := s.repository.Add(ctx, comment)

	if err != nil {
		return 0, err
	}

	return upComment.ID, nil
}
