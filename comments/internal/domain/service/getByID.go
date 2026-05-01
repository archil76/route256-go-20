package comments

import (
	"context"
	"route256/comments/internal/domain/model"
)

func (s *CommentsService) GetByID(ctx context.Context, ID int64) (*model.Comment, error) {

	upComment, err := s.repository.GetByID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return upComment, nil
}
