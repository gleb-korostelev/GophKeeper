package db

import "time"

// Config defines the configuration settings required for initializing
// and managing a database connection pool.
type Config struct {
	MaxOpenConns    int           // Maximum number of open connections in the pool.
	MaxIdleConns    int           // Maximum number of idle connections in the pool.
	ConnMaxLifetime time.Duration // Maximum lifetime of a connection in the pool.
	Dsn             string        // Data Source Name for the database connection.
}
