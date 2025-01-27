package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
)

// mock is a test implementation of the IAdapter interface.
//
// This mock can be used to simulate database behavior in unit tests, enabling
// testing of logic that depends on database operations without requiring a live database.
type mock struct{}

// GetConn returns a nil *sql.DB, simulating a disconnected state.
func (m *mock) GetConn(context.Context) *sql.DB {
	return nil
}

// InTx simulates a transaction execution by calling the provided function
// without performing any actual database operations.
func (m *mock) InTx(context.Context, func(ctx context.Context, tx pgx.Tx) error) error {
	return nil
}
