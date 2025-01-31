// Package config provides configuration keys for retrieving application settings.
// These keys are used to fetch environment variables or other configuration values
// required by the application.
package config

type configKey string

const (
	// Port specifies the port on which the HTTP server will run.
	Port = configKey("PORT")

	// DBDSN specifies the Data Source Name (DSN) for the database connection.
	DBDSN = configKey("DB_DSN")

	// HttpsHost specifies the hostname for the HTTPS server.
	HttpsHost = configKey("HTTPS_HOST")

	// IsSwaggerCreated determines whether Swagger documentation should be created and exposed.
	IsSwaggerCreated = configKey("IS_SWAGGER_CREATED")

	// JwtKey specifies the private key for signing JWT tokens.
	JwtKey = configKey("JWT_KEY")

	// MaxOpenConns specifies the maximum number of open database connections.
	MaxOpenConns = configKey("MAX_OPEN_CONNS")

	// MaxIdleConns specifies the maximum number of idle database connections.
	MaxIdleConns = configKey("MAX_IDLE_CONNS")

	// ConnMaxLifetime specifies the maximum lifetime of a database connection.
	ConnMaxLifetime = configKey("CONN_MAX_LIFETIME")
)
