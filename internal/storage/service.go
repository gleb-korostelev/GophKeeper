package storage

import "context"

type Storage interface {
	Ping(ctx context.Context) (int, error)
	Close() error
}
