// Package profile provides services for managing user card information in the GophKeeper application.
package profile

import (
	"context"
	"fmt"

	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/gleb-korostelev/GophKeeper/repository"
	"github.com/gleb-korostelev/GophKeeper/tools/db"
	"github.com/jackc/pgx/v5"
)

// service defines the implementation of the profile service.
//
// Fields:
// - db: The database adapter for executing transactional operations.
type service struct {
	db   db.IAdapter
	repo repository.Repository
}

// NewService creates a new instance of the profile service.
func NewService(db db.IAdapter) *service {
	return &service{db: db}
}

// UploadInfo uploads or updates a user's card information in the database.
func (s *service) UploadInfo(ctx context.Context, profile profile.CardInfo) (err error) {
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err = s.repo.UploadCardInfo(ctx, tx, profile)
		if err != nil {
			return fmt.Errorf("error in uploadCardInfo: %w", err)
		}
		return nil
	})
	return
}

// GetUserCards retrieves all card information associated with a username.
func (s *service) GetUserCards(ctx context.Context, username string) (profile []profile.CardInfo, err error) {
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		profile, err = s.repo.GetUserCards(ctx, tx, username)
		if err != nil {
			return fmt.Errorf("error in getUserCards: %w", err)
		}
		return nil
	})
	return
}

// DeleteCard deletes a specific card associated with a username.
func (s *service) DeleteCard(ctx context.Context, username, cardNumber string) (err error) {
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err = s.repo.DeleteCard(ctx, tx, username, cardNumber)
		if err != nil {
			return fmt.Errorf("error in deleteCard: %w", err)
		}
		return nil
	})
	return
}
