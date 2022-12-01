package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type contextKey struct{}

var contextKeyTransaction contextKey

func ctxWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, contextKeyTransaction, tx)
}

func txFromCtx(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(contextKeyTransaction).(pgx.Tx)
	if ok {
		return tx
	}
	return nil
}
