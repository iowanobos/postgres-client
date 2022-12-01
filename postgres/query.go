package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryManager interface {
	Select(ctx context.Context, builder squirrel.SelectBuilder) (pgx.Rows, error)
	Insert(ctx context.Context, builder squirrel.InsertBuilder) (int64, error)
	Update(ctx context.Context, builder squirrel.UpdateBuilder) (int64, error)
	Delete(ctx context.Context, builder squirrel.DeleteBuilder) (int64, error)
	Execute(ctx context.Context, query string, args ...any) (int64, error)
}

type queryManager struct {
	pool *pgxpool.Pool
}

func (c *queryManager) Select(ctx context.Context, builder squirrel.SelectBuilder) (pgx.Rows, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	if tx := txFromCtx(ctx); tx != nil {
		return tx.Query(ctx, query, args...)
	}
	return c.pool.Query(ctx, query, args...)
}

func (c *queryManager) Insert(ctx context.Context, builder squirrel.InsertBuilder) (int64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	return c.Execute(ctx, query, args...)
}

func (c *queryManager) Update(ctx context.Context, builder squirrel.UpdateBuilder) (int64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	return c.Execute(ctx, query, args...)
}

func (c *queryManager) Delete(ctx context.Context, builder squirrel.DeleteBuilder) (int64, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	return c.Execute(ctx, query, args...)
}

func (c *queryManager) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	if tx := txFromCtx(ctx); tx != nil {
		return tx.Query(ctx, query, args...)
	}
	return c.pool.Query(ctx, query, args...)
}

func (c *queryManager) Execute(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := c.execute(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func (c *queryManager) execute(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	if tx := txFromCtx(ctx); tx != nil {
		return tx.Exec(ctx, query, args...)
	}
	return c.pool.Exec(ctx, query, args...)
}
