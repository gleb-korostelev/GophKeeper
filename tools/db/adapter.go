// Package db provides a database adapter for interacting with PostgreSQL using the pgx library.
// It supports connection pooling, transactional operations, and database migrations with Goose.
//
// This package is designed to simplify database interactions for the GophKeeper application.

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	gophkeeper "github.com/gleb-korostelev/GophKeeper"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// IAdapter defines the interface for the database adapter.
//
// Methods:
// - InTx: Executes a function within a database transaction.
// - GetConn: Returns the underlying connection pool.
type IAdapter interface {
	InTx(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error
	GetConn() *pgxpool.Pool
}

// Adapter implements the IAdapter interface and provides database connectivity
// and transactional support.
type Adapter struct {
	pool      *pgxpool.Pool      // Connection pool for the database.
	isolation sql.IsolationLevel // Isolation level for transactions.
	Config    *Config            // Configuration for the adapter.
}

// NewAdapter initializes a new database adapter.
func NewAdapter(ctx context.Context, config Config, isolation sql.IsolationLevel) (IAdapter, error) {
	poolConfig, err := pgxpool.ParseConfig(config.Dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	poolConfig.MaxConns = int32(config.MaxOpenConns)
	poolConfig.MaxConnLifetime = config.ConnMaxLifetime * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	ad := &Adapter{
		pool:      pool,
		isolation: isolation,
		Config:    &config,
	}

	if err := ad.gooseUp(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}
	return ad, nil
}

// GetConn returns the underlying connection pool.
func (b *Adapter) GetConn() *pgxpool.Pool {
	return b.pool
}

// InTx executes a function within a database transaction.
func (b *Adapter) InTx(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) (err error) {
	tx, err := b.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error creating tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			logger.Error(p)
			err = fmt.Errorf("panic: %v", p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = f(ctx, tx)
	return
}

// gooseUp applies database migrations using Goose.
func (b *Adapter) gooseUp(ctx context.Context, pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer sqlDB.Close()
	goose.SetBaseFS(gophkeeper.EmbedMigrations)

	if err := goose.Up(sqlDB, "migrations", goose.WithAllowMissing()); err != nil {
		return fmt.Errorf("goose up error: %w", err)
	}
	return nil
}

// gooseCreate creates a new migration file using Goose.
func (b *Adapter) gooseCreate(ctx context.Context, pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer sqlDB.Close()
	goose.SetBaseFS(gophkeeper.EmbedMigrations)
	if err := goose.Create(sqlDB, "migrations", "", "sql"); err != nil {
		return err
	}
	return nil
}

// gooseDown rolls back the last database migration using Goose.
func (b *Adapter) gooseDown(ctx context.Context, pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer sqlDB.Close()
	goose.SetBaseFS(gophkeeper.EmbedMigrations)
	if err := goose.Down(sqlDB, "migrations"); err != nil {
		return err
	}

	return nil
}
