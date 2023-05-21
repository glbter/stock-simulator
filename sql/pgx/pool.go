package pgx

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

func NewPool(ctx context.Context, dsn string) (*Pool, error) {
	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	return &Pool{pool: pool}, nil
}

func (p *Pool) Close() error {
	p.pool.Close()
	return nil
}
