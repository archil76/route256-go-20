package postgres

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (r *Repository) Reserve(ctx context.Context, items []model.Item) ([]model.Stock, error) {
	const query = `SELECT id, total_count, reserved FROM stocks WHERE id = ANY($1)`

	listOfID := make([]int64, len(items))
	itemsMap := map[int64]*model.Item{}
	for _, item := range items {
		listOfID = append(listOfID, item.Sku)
		itemsMap[item.Sku] = &item
	}

	rows, err := r.pool.Query(ctx, query, listOfID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []model.Stock
	for rows.Next() {
		upStock := &model.Stock{}
		if err := rows.Scan(&upStock.Sku, &upStock.TotalCount, &upStock.Reserved); err != nil {
			return nil, err
		}

		item := itemsMap[upStock.Sku]
		available := upStock.TotalCount - upStock.Reserved
		if available < item.Count {
			return nil, model.ErrShortOfStock
		}

		stocks = append(stocks, model.Stock{
			Sku:        upStock.Sku,
			TotalCount: upStock.TotalCount,
			Reserved:   upStock.Reserved + item.Count,
		})
	}

	for _, stock := range stocks {
		_, err := r.updateStock(ctx, stock)
		if err != nil {
			return nil, err
		}
	}

	return stocks, nil
}
