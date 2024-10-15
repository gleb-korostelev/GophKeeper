package db

import (
	"context"
	"database/sql"
	"fmt"

	. "github.com/gleb-korostelev/GophKeeper"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type IAdapter interface {
	InTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error
	GetConn(ctx context.Context) *sql.DB
}

type Adapter struct {
	Conn      *sql.DB
	isolation sql.IsolationLevel
}

func NewAdapter(conn *sql.DB, isolation sql.IsolationLevel) (IAdapter, error) {
	ad := &Adapter{Conn: conn, isolation: isolation}
	err := ad.GooseUp()
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (b *Adapter) GetConn(ctx context.Context) *sql.DB {
	return b.Conn
}

func (b *Adapter) InTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) (err error) {
	tx, err := b.Conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: b.isolation,
	})
	if err != nil {
		return fmt.Errorf("error creating tx: %s", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			logger.Error(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = f(ctx, tx)
	return
}

func (b *Adapter) GooseUp() error {
	goose.SetBaseFS(EmbedMigrations)
	if err := goose.Up(b.Conn, "migrations", goose.WithAllowMissing()); err != nil {
		return err
	}
	return nil
}

func (b *Adapter) GooseCreate() error {
	goose.SetBaseFS(EmbedMigrations)
	if err := goose.Create(b.Conn, "migrations", "", "sql"); err != nil {
		return err
	}
	return nil
}

func (b *Adapter) GooseDown() error {
	goose.SetBaseFS(EmbedMigrations)
	if err := goose.Down(b.Conn, "migrations"); err != nil {
		return err
	}

	return nil
}
