// Package initConnection provides functions to initialize core application components,
// including HTTP handlers, services, and middleware for GophKeeper.
package initConnection

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/gleb-korostelev/GophKeeper/config"
	"github.com/gleb-korostelev/GophKeeper/tools/db"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
)

// NewDBConn initializes a new database connection adapter using the provided context
// and configuration parameters from the application's configuration.
//
// It retrieves the following configuration settings:
//   - config.DBDSN: the database connection string.
//   - config.MaxOpenConns: the maximum number of open database connections.
//   - config.MaxIdleConns: the maximum number of idle connections.
//   - config.ConnMaxLifetime: the maximum lifetime of a connection.
func NewDBConn(ctx context.Context) db.IAdapter {
	dsn := config.GetConfigString(config.DBDSN)

	cfg := db.Config{
		MaxOpenConns:    config.GetConfigInt(config.MaxOpenConns),
		MaxIdleConns:    config.GetConfigInt(config.MaxIdleConns),
		ConnMaxLifetime: config.GetConfigDuration(config.ConnMaxLifetime),
		Dsn:             dsn,
	}

	ad, err := db.NewAdapter(ctx, cfg, sql.LevelReadUncommitted)
	if err != nil {
		logger.Fatal("can't initialize database adapter: %w", err.Error())
	}
	return ad
}
