package apperror

import "errors"

var (
	ErrNoServerAddress       = errors.New("server address is empty")
	ErrNoDatabaseDestination = errors.New("database destination is empty")
)
