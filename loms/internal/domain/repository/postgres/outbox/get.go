package postgres

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) Get(ctx context.Context) (*[]model.OutboxItem, error) {
	pool := r.pooler.PickPool(ctx)

	const query = `select id, key, payload from outbox where status = 'new'`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	Items := []model.OutboxItem{}
	for rows.Next() {
		item := model.OutboxItem{}
		if err := rows.Scan(&item.Id, &item.Key, &item.Payload); err != nil {
			return nil, err
		}
		Items = append(Items, model.OutboxItem{
			Id:      item.Id,
			Key:     item.Key,
			Payload: item.Payload,
		})
	}

	return &Items, nil
}
