package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

type Options struct {
	ConnString string
}

func New(ctx context.Context, options Options) (*Client, error) {
	config, err := pgxpool.ParseConfig(options.ConnString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Client{
		pool: pool,
	}, nil
}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) TransactionManager() TransactionManager {
	return &transactionManager{pool: c.pool}
}

func (c *Client) QueryManager() QueryManager {
	return &queryManager{pool: c.pool}
}
