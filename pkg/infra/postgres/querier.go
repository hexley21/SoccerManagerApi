package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func Rollback(ctx context.Context, tx pgx.Tx, internal error) error {
	if err := tx.Rollback(ctx); err != nil {
		return errors.Join(
			fmt.Errorf("rollback failed: %w", err),
			fmt.Errorf("internal: %w", internal),
		)
	}
	return internal
}
