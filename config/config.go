// Package config provides configuration keys for retrieving application settings.
// These keys are used to fetch environment variables or other configuration values
// required by the application.
package config

// configKey defines a type for configuration keys.
// This type is used to identify specific application settings in a strongly-typed manner.
type configKey string

const (
	// Port specifies the port on which the HTTP server will run.
	// Example: "8080"
	Port = configKey("PORT")

	// DBDSN specifies the Data Source Name (DSN) for the database connection.
	// Example: "postgres://user:password@localhost:5432/dbname"
	DBDSN = configKey("DB_DSN")

	// HttpsHost specifies the hostname for the HTTPS server.
	// Example: "example.com"
	HttpsHost = configKey("HTTPS_HOST")

	// IsSwaggerCreated determines whether Swagger documentation should be created and exposed.
	// Example: "true" or "false"
	IsSwaggerCreated = configKey("IS_SWAGGER_CREATED")

	// JwtKey specifies the private key for signing JWT tokens.
	// This is expected to be a hex-encoded string.
	// Example: "a1b2c3d4e5..."
	JwtKey = configKey("JWT_KEY")

	// MaxOpenConns specifies the maximum number of open database connections.
	// Example: "50"
	MaxOpenConns = configKey("MAX_OPEN_CONNS")

	// MaxIdleConns specifies the maximum number of idle database connections.
	// Example: "10"
	MaxIdleConns = configKey("MAX_IDLE_CONNS")

	// ConnMaxLifetime specifies the maximum lifetime of a database connection.
	// This value is typically represented as a duration string.
	// Example: "30m"
	ConnMaxLifetime = configKey("CONN_MAX_LIFETIME")
)
