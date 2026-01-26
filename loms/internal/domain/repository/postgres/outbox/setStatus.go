package postgres

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
)

func (r *Repository) SetStatus(ctx context.Context, id int, status string) error {
	var err error

	err = r.pooler.InTx(ctx, func(ctx context.Context) error {
		pool := r.pooler.PickPool(ctx)

		const query = `UPDATE outbox SET status=$2, sent_at=$3 where id = $1`
		if _, err := pool.Exec(ctx, query, id, status, time.Now()); err != nil {
			return errors.Wrap(err, "pgx.QueryRow.Scan")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
