package postgres

import (
	"context"

	"github.com/cockroachdb/errors"
)

func (r *Repository) Create(ctx context.Context, key, status string, payload []byte) (int, error) {
	var err error
	id := 0
	err = r.pooler.InTx(ctx, func(ctx context.Context) error {
		pool := r.pooler.PickPool(ctx)

		const queryOrders = `INSERT INTO outbox (key, payload, status) VALUES ($1, $2, $3) returning id`
		if err = pool.QueryRow(ctx, queryOrders, key, payload, status).
			Scan(&id); err != nil {
			err = errors.Wrap(err, "pgx.QueryRow.Scan")
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
