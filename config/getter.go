// Package config provides utilities for accessing application configuration values
// from environment variables. It includes helper functions to retrieve configuration
// values in various data types, ensuring validation and error handling.
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/gleb-korostelev/GophKeeper/tools/logger"
)

// tempErr is a format string used for error messages when a configuration parameter is invalid.
const tempErr = "Config parameter %s is incorrect: %s"

// GetConfigString retrieves the value of an environment variable as a string.
// If the environment variable is not set or is empty, the application will terminate with a fatal error.
//
// Parameters:
//   - key: The configuration key to retrieve.
//
// Returns:
//   - The value of the environment variable as a string.
//
// Example usage:
//
//	dbDSN := config.GetConfigString(config.DBDSN)
func GetConfigString(key configKey) string {
	value := os.Getenv(string(key))
	if value == "" {
		logger.Fatalf(tempErr, key, value)
	}
	return value
}

// GetConfigBool retrieves the value of an environment variable as a boolean.
// If the value cannot be parsed as a boolean, the application will terminate with a fatal error.
//
// Parameters:
//   - key: The configuration key to retrieve.
//
// Returns:
//   - The value of the environment variable as a boolean.
//
// Example usage:
//
//	isSwaggerEnabled := config.GetConfigBool(config.IsSwaggerCreated)
func GetConfigBool(key configKey) bool {
	valueStr := os.Getenv(string(key))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, err)
	}
	return value
}

// GetConfigInt retrieves the value of an environment variable as an integer.
// If the value cannot be parsed as an integer, the application will terminate with a fatal error.
//
// Parameters:
//   - key: The configuration key to retrieve.
//
// Returns:
//   - The value of the environment variable as an integer.
//
// Example usage:
//
//	maxOpenConns := config.GetConfigInt(config.MaxOpenConns)
func GetConfigInt(key configKey) int {
	valueStr := os.Getenv(string(key))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, valueStr)
	}
	return value
}

// GetConfigDuration retrieves the value of an environment variable as a time.Duration.
// If the value cannot be parsed as a valid duration string, the application will terminate with a fatal error.
//
// Parameters:
//   - key: The configuration key to retrieve.
//
// Returns:
//   - The value of the environment variable as a time.Duration.
//
// Example usage:
//
//	connLifetime := config.GetConfigDuration(config.ConnMaxLifetime)
func GetConfigDuration(key configKey) time.Duration {
	valueStr := os.Getenv(string(key))
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, valueStr)
	}
	return value
}
