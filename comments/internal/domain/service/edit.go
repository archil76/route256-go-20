package comments

import (
	"context"
	"time"

	"route256/comments/internal/domain/model"
)

func (s *CommentsService) Edit(ctx context.Context, comment model.Comment) (int64, error) {

	exComment, err := s.repository.GetByID(ctx, comment.ID)
	if err != nil {
		return 0, model.ErrCommentDoesntExist
	}

	// Проверяем, что редактирующий пользователь является автором комментария
	if exComment.UserID != comment.UserID {
		return 0, model.ErrUserNotAuthor
	}

	// Проверяем, что не прошло больше editInterval секунд с момента создания
	if time.Since(exComment.CreatedAt) > s.editInterval*time.Second {
		return 0, model.ErrEditTimeExpired
	}

	upComment, err := s.repository.Edit(ctx, comment)
	if err != nil {
		return 0, err
	}

	return upComment.ID, nil
}
