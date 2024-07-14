package pgxtx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryInTx func(context.Context, pgx.Tx) error

func InTx(ctx context.Context, pool *pgxpool.Pool, fns ...QueryInTx) (err error) {
	var tx pgx.Tx
	tx, err = pool.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin transaction: %w", err)
		return
	}

	defer func() {
		rollback := func() {
			if rErr := tx.Rollback(ctx); rErr != nil {
				l.Errorf("transaction rollback failed: %s", rErr.Error())
			}
		}
		if p := recover(); p != nil {
			l.Errorf("panic recovered in transaction: %v", p)
			rollback()
			panic(p)
		}

		if err != nil {
			rollback()
			return
		}

		err = tx.Commit(ctx)
	}()

	err = executeQueries(ctx, tx, fns)
	return
}

func executeQueries(ctx context.Context, tx pgx.Tx, queries []QueryInTx) error {
	for _, q := range queries {
		if err := q(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
