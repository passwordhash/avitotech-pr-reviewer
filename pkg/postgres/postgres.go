package pgPkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)


var ErrNewPool = errors.New("failed to create new pgx pool")

type Option func(*pgxpool.Config)

func WithMaxConns(maxConns int32) Option {
	return func(cfg *pgxpool.Config) {
		if maxConns > 0 {
			cfg.MaxConns = maxConns
		}
	}
}

// NewPool создает новый пул подключений к базе данных PostgreSQL.
func NewPool(ctx context.Context, dsn string, opts ...Option) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNewPool, err.Error())
	}

	for _, opt := range opts {
		opt(config)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrNewPool, err.Error())
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()

		return nil, fmt.Errorf("%w: %s", ErrNewPool, err.Error())
	}

	return pool, nil
}

type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Batch = pgx.Batch

// type Queryer interface {
//	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
//}
