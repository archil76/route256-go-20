package postgres

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
)

type stockItem struct {
	Count      uint32
	TotalCount uint32
	Reserved   uint32
}

func (r *Repository) Reserve(ctx context.Context, items []model.Item) ([]model.Stock, error) {
	var err error
	var stocks []model.Stock

	pool := r.pooler.PickPool(ctx)

	const query = `SELECT id, total_count, reserved FROM stocks WHERE id = ANY($1::bigint[]) FOR UPDATE`

	//listOfID := make([]string, len(items))
	listOfID := make([]int64, len(items))
	itemsMap := map[int64]*stockItem{}
	for i, item := range items {
		//listOfID[i] = fmt.Sprintf("%d", item.Sku)
		listOfID[i] = item.Sku
		itemsMap[item.Sku] = &stockItem{
			Count: item.Count,
		}
	}
	var rows pgx.Rows
	//rows, err = pool.Query(ctx, fmt.Sprintf(query, strings.Join(listOfID, ",")))
	rows, err = pool.Query(ctx, query, listOfID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		upStock := &model.Stock{}
		if err = rows.Scan(&upStock.Sku, &upStock.TotalCount, &upStock.Reserved); err != nil {
			return nil, err
		}

		item := itemsMap[upStock.Sku]
		item.TotalCount = upStock.TotalCount
		item.Reserved = upStock.Reserved
	}

	for sku, item := range itemsMap {
		available := item.TotalCount - item.Reserved
		if available < item.Count {
			err = model.ErrOutOfStock
			return nil, err
		}

		stocks = append(stocks, model.Stock{
			Sku:        sku,
			TotalCount: item.TotalCount,
			Reserved:   item.Reserved + item.Count,
		})
	}

	const queryUpdateStock = `UPDATE stocks SET reserved=$2 where id = $1`
	batch := &pgx.Batch{}
	for _, stock := range stocks {
		batch.Queue(queryUpdateStock, stock.Sku, stock.Reserved)
	}
	br := pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err = br.Exec()
		if err != nil {
			err = errors.Wrap(err, "pgx.QueryRow.Scan")
			return nil, err
		}
	}

	return stocks, nil
}
