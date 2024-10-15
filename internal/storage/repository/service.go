package repository

import (
	"context"
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/internal/db"
	"github.com/gleb-korostelev/GophKeeper/internal/storage"
	logger "github.com/gleb-korostelev/GophKeeper/tools"
)

type service struct {
	data db.DB
}

func NewDBStorage(data db.DB) storage.Storage {
	return &service{
		data: data,
	}
}

func (s *service) Ping(ctx context.Context) (int, error) {
	err := s.data.Ping(context.Background())
	if err != nil {
		logger.Errorf("Failed to connect to the database %v", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *service) Close() error {
	err := s.data.Close()
	if err != nil {
		return err
	}
	return nil
}
