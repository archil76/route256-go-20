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
	shardIndex := r.sm.GetShardIndex(sku)
	pool, err := r.sm.PickPool(ctx, shardIndex)
	if err != nil {
		return nil, err
	}

	currentList, err := r.getList(ctx, pool, query, sku)
	if err != nil {
		return nil, err
	}

	sortCommentsList(currentList.Comments)

	return &model.CommentsList{Comments: currentList.Comments}, nil
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

	sortCommentsList(allComments)

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

func sortCommentsList(currentList []model.Comment) {
	if len(currentList) != 0 {
		sort.Slice(currentList, func(i, j int) bool {
			// Сначала сортируем по CreatedAt по убыванию (новые первыми)
			if currentList[i].CreatedAt.Equal(currentList[j].CreatedAt) {
				// Если CreatedAt равны, сортируем по UserID по возрастанию
				return currentList[i].UserID < currentList[j].UserID
			}
			return currentList[i].CreatedAt.After(currentList[j].CreatedAt)
		})
	}
}
