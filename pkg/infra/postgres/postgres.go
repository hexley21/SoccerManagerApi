package postgres

import (
	"context"
	"fmt"

	"github.com/hexley21/soccer-manager/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnUrl(cfg config.Postgres) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SslMode,
	)
}

func New(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(GetConnUrl(cfg))
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	pgPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pgPool.Ping(ctx); err != nil {
		return nil, err
	}

	return pgPool, nil
}

func MigrateCfg(cfg config.Postgres, migrationPath string) (*migrate.Migrate, error) {
	return Migrate(GetConnUrl(cfg), migrationPath)
}

func Migrate(url string, migrationPath string) (*migrate.Migrate, error) {
	m, err := migrate.New(migrationPath, url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect migrator: %w", err)
	}

	return m, nil
}
