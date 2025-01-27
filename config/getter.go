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
func GetConfigString(key configKey) string {
	value := os.Getenv(string(key))
	if value == "" {
		logger.Fatalf(tempErr, key, value)
	}
	return value
}

// GetConfigBool retrieves the value of an environment variable as a boolean.
func GetConfigBool(key configKey) bool {
	valueStr := os.Getenv(string(key))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, err)
	}
	return value
}

// GetConfigInt retrieves the value of an environment variable as an integer.
func GetConfigInt(key configKey) int {
	valueStr := os.Getenv(string(key))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, valueStr)
	}
	return value
}

// GetConfigDuration retrieves the value of an environment variable as a time.Duration.
func GetConfigDuration(key configKey) time.Duration {
	valueStr := os.Getenv(string(key))
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		logger.Fatalf(tempErr, key, valueStr)
	}
	return value
}
