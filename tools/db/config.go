package db

import "time"

// Config defines the configuration settings required for initializing
// and managing a database connection pool.
//
// Fields:
// - MaxOpenConns: Maximum number of open connections allowed in the pool. This limits
//   the number of simultaneous database queries and helps prevent overloading the database server.
// - MaxIdleConns: Maximum number of idle connections maintained in the pool. This allows
//   faster execution of queries by reusing idle connections.
// - ConnMaxLifetime: The maximum duration a connection can remain open before being closed and replaced.
//   Helps prevent issues caused by stale connections.
// - Dsn: The Data Source Name (DSN) for the database connection. It contains the connection string
//   required to connect to the database.
//
// Example Usage:
//   config := db.Config{
//       MaxOpenConns:    10,
//       MaxIdleConns:    5,
//       ConnMaxLifetime: 30 * time.Minute,
//       Dsn:             "postgresql://user:password@localhost/dbname",
//   }
type Config struct {
	MaxOpenConns    int           // Maximum number of open connections in the pool.
	MaxIdleConns    int           // Maximum number of idle connections in the pool.
	ConnMaxLifetime time.Duration // Maximum lifetime of a connection in the pool.
	Dsn             string        // Data Source Name for the database connection.
}
