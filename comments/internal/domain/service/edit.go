package comments

import (
	"context"
	"route256/comments/internal/domain/model"
)

func (s *CommentsService) Edit(ctx context.Context, comment model.Comment) (int64, error) {

	upComment, err := s.repository.Edit(ctx, comment)

	if err != nil {
		return 0, err
	}

	return upComment.ID, nil
}
