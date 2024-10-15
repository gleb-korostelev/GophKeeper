package db

import (
	"context"
	"database/sql"
)

type mock struct{}

func NewDBMock() IAdapter {
	return &mock{}
}

func (m *mock) GetConn(context.Context) *sql.DB {
	return nil
}

func (m *mock) InTx(context.Context, func(ctx context.Context, tx *sql.Tx) error) error {
	return nil
}
