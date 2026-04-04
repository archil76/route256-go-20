package postgres

import (
	"context"
	"route256/comments/internal/domain/model"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

// узнаем шард и получаем список
func (r *Repository) GetListBySKU(ctx context.Context, sku int64) (*model.CommentsList, error) {

	if sku < 1 {
		return nil, model.ErrSkuIsNotValid
	}
	const query = `SELECT id, user_id, sku, comment, created_at FROM comments where sku = $1`
	pool, err := r.sm.PickPool(ctx, sku)
	if err != nil {
		return nil, err
	}

	return r.getList(ctx, pool, query, sku)
}

// ходим по всем шардам, получаем списки и соединяем их
func (r *Repository) GetListByUserID(ctx context.Context, userID int64) (*model.CommentsList, error) {

	if userID < 1 {
		return nil, model.ErrUserIDIsNotValid
	}
	const query = `SELECT id, user_id, sku, comment, created_at FROM comments where user_id = $1`

	// Собираем все комментарии со всех шардов
	var allComments []model.Comment

	for _, pool := range r.sm.GetShards(ctx) {

		currentList, err := r.getList(ctx, pool, query, userID)
		if err != nil {
			return nil, err
		}

		if len(currentList.Comments) != 0 {
			allComments = append(allComments, currentList.Comments...)
		}
	}

	if len(allComments) != 0 {
		sort.Slice(allComments, func(i, j int) bool {
			return allComments[i].CreatedAt.Before(allComments[j].CreatedAt)
		})
	}

	return &model.CommentsList{Comments: allComments}, nil
}

func (r *Repository) getList(ctx context.Context, pool *pgxpool.Pool, query string, conditionId int64) (*model.CommentsList, error) {

	rows, err := pool.Query(ctx, query, conditionId)
	if err != nil {
		return nil, err
	}

	commentList := model.CommentsList{}
	for rows.Next() {
		upComment := model.Comment{}
		if err := rows.Scan(&upComment.ID, &upComment.UserID, &upComment.Sku, &upComment.Comment, &upComment.CreatedAt); err != nil {
			return nil, err
		}
		commentList.Comments = append(commentList.Comments, model.Comment{
			ID:        upComment.ID,
			UserID:    upComment.UserID,
			Sku:       upComment.Sku,
			Comment:   upComment.Comment,
			CreatedAt: upComment.CreatedAt,
		})
	}

	return &commentList, nil
}
