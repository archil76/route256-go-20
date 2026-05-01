package comments

import (
	"context"
	"route256/comments/internal/domain/model"
)

func (s *CommentsService) GetListBySKU(ctx context.Context, sku int64) (*model.CommentsList, error) {

	upCommentsList, err := s.repository.GetListBySKU(ctx, sku)

	if err != nil {
		return nil, err
	}

	return upCommentsList, nil
}
