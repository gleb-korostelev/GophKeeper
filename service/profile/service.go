// Package profile provides services for managing user card information in the GophKeeper application.
package profile

import (
	"context"
	"fmt"

	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/gleb-korostelev/GophKeeper/tools/db"
	"github.com/jackc/pgx/v5"
)

// service defines the implementation of the profile service.
//
// Fields:
// - db: The database adapter for executing transactional operations.
type service struct {
	db db.IAdapter
}

// NewService creates a new instance of the profile service.
//
// Parameters:
// - db: The database adapter for database operations.
//
// Returns:
// - *service: The profile service implementation.
func NewService(db db.IAdapter) *service {
	return &service{db: db}
}

// UploadInfo uploads or updates a user's card information in the database.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - profile: The `profile.CardInfo` containing the card details to upload or update.
//
// Returns:
// - err: An error if the operation fails.
//
// Workflow:
// 1. Executes a database transaction.
// 2. Calls `uploadCardInfo` to insert or update the card information.
// 3. Returns any error encountered during the process.
func (s *service) UploadInfo(ctx context.Context, profile profile.CardInfo) (err error) {
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err = uploadCardInfo(ctx, tx, profile)
		if err != nil {
			return fmt.Errorf("error in uploadCardInfo: %w", err)
		}
		return nil
	})
	return
}

// GetUserCards retrieves all card information associated with a username.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - username: The username for which card information is retrieved.
//
// Returns:
// - profile: A slice of `profile.CardInfo` containing the user's card details.
// - err: An error if the operation fails.
//
// Workflow:
// 1. Executes a database transaction.
// 2. Calls `getUserCards` to retrieve card details from the database.
// 3. Returns the card details and any error encountered.
func (s *service) GetUserCards(ctx context.Context, username string) (profile []profile.CardInfo, err error) {
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		profile, err = getUserCards(ctx, tx, username)
		if err != nil {
			return fmt.Errorf("error in getUserCards: %w", err)
		}
		return nil
	})
	return
}

// DeleteCard deletes a specific card associated with a username.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - username: The username associated with the card.
// - cardNumber: The card number to delete.
//
// Returns:
// - err: An error if the deletion fails or the card does not exist.
//
// Workflow:
// 1. Executes a database transaction.
// 2. Calls `deleteCard` to remove the card information from the database.
// 3. Returns any error encountered during the process.
func (s *service) DeleteCard(ctx context.Context, username, cardNumber string) (err error) {
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err = deleteCard(ctx, tx, username, cardNumber)
		if err != nil {
			return fmt.Errorf("error in deleteCard: %w", err)
		}
		return nil
	})
	return
}
