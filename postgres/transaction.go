package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionManager handles transactions
type TransactionManager interface {
	BeginTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type transactionManager struct {
	pool *pgxpool.Pool
}

func (c *transactionManager) BeginTx(ctx context.Context) (context.Context, error) {
	tx := txFromCtx(ctx)
	if tx == nil { // TODO: Пока что недопустимы вложенные транзакции
		var err error
		tx, err = c.pool.Begin(ctx)
		if err != nil {
			return nil, err
		}

		return ctxWithTx(ctx, tx), nil
	}
	return ctx, nil

}

func (c *transactionManager) CommitTx(ctx context.Context) error {
	tx := txFromCtx(ctx)
	if tx == nil {
		return nil // TODO: если нет транзакции в контексте то ничего не делаю?
	}

	return tx.Commit(ctx)
}

func (c *transactionManager) Rollback(ctx context.Context) error {
	tx := txFromCtx(ctx)
	if tx == nil {
		return nil // TODO: если нет транзакции в контексте то ничего не делаю?
	}

	return tx.Rollback(ctx)
}
