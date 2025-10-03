package postgres

import (
	"context"
	"fmt"
	"route256/loms/internal/domain/model"
	"strings"
)

type stockItem struct {
	Count      uint32
	TotalCount uint32
	Reserved   uint32
}

func (r *Repository) Reserve(ctx context.Context, items []model.Item) ([]model.Stock, error) {
	const query = `SELECT id, total_count, reserved FROM stocks WHERE id in (%s)`

	listOfID := make([]string, len(items))
	itemsMap := map[int64]*stockItem{}
	for i, item := range items {
		listOfID[i] = fmt.Sprintf("%d", item.Sku)
		itemsMap[item.Sku] = &stockItem{
			Count: item.Count,
		}
	}

	rows, err := r.pool.Query(ctx, fmt.Sprintf(query, strings.Join(listOfID, ",")))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		upStock := &model.Stock{}
		if err := rows.Scan(&upStock.Sku, &upStock.TotalCount, &upStock.Reserved); err != nil {
			return nil, err
		}

		item := itemsMap[upStock.Sku]
		item.TotalCount = upStock.TotalCount
		item.Reserved = upStock.Reserved
	}
	var stocks []model.Stock
	for sku, item := range itemsMap {
		available := item.TotalCount - item.Reserved
		if available < item.Count {
			return nil, model.ErrShortOfStock
		}

		stocks = append(stocks, model.Stock{
			Sku:        sku,
			TotalCount: item.TotalCount,
			Reserved:   item.Reserved + item.Count,
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
