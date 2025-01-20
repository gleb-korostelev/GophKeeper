package profile

import (
	"context"
	"fmt"

	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/jackc/pgx/v5"
)

// getAccountByUserName retrieves account information by username.
//
// Parameters:
// - ctx: Context for managing request deadlines and cancellations.
// - tx: Database transaction interface.
// - username: The username of the account to retrieve.
//
// Returns:
// - acc: The retrieved account information as a `models.Account`.
// - err: An error if the query or scan fails.
func getAccountByUserName(ctx context.Context, tx pgx.Tx, username string) (acc models.Account, err error) {
	const query = `
		SELECT id,
			   username,
			   secret,
			   account_type,
			   created_at,
			   role_changed_at,
			   updated_at
		FROM auth.users
		WHERE username = $1;
	`

	var user string
	err = tx.QueryRow(ctx, query, username).Scan(
		&acc.ID,
		&user,
		&acc.Secret,
		&acc.AccountType,
		&acc.CreatedAt,
		&acc.RoleChangedAt,
		&acc.UpdatedAt,
	)
	acc.Username = user
	return
}

// uploadCardInfo uploads or updates card information for a user.
//
// Parameters:
// - ctx: Context for managing request deadlines and cancellations.
// - tx: Database transaction interface.
// - profile: Card information to upload or update.
//
// Returns:
// - err: An error if the operation fails.
func uploadCardInfo(ctx context.Context, tx pgx.Tx, profile profile.CardInfo) (err error) {
	const query = `
    INSERT INTO auth.cards (user_id, card_holder, card_number, expiration_date, cvv, metadata, updated_at)
    SELECT id, $2, $3, $4, $5, $6, now()
    FROM auth.users
    WHERE username = $1
    ON CONFLICT (user_id, card_number)
    DO UPDATE SET 
        card_holder = EXCLUDED.card_holder,
        expiration_date = EXCLUDED.expiration_date,
        cvv = EXCLUDED.cvv,
        metadata = EXCLUDED.metadata,
        updated_at = now();
    `

	_, err = tx.Exec(ctx, query,
		profile.Username,
		profile.CardHolder,
		profile.CardNumber,
		profile.ExpirationDate,
		profile.Cvv,
		profile.Metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to upload card info: %w", err)
	}

	return
}

// getUserCards retrieves all cards associated with a user.
//
// Parameters:
// - ctx: Context for managing request deadlines and cancellations.
// - tx: Database transaction interface.
// - username: The username of the user whose cards are being retrieved.
//
// Returns:
// - cards: A slice of `profile.CardInfo` containing card details.
// - err: An error if the query or scanning fails.
func getUserCards(ctx context.Context, tx pgx.Tx, username string) ([]profile.CardInfo, error) {
	var cards []profile.CardInfo

	const query = `
        SELECT c.card_number, c.card_holder, c.expiration_date, c.cvv, c.metadata
        FROM auth.cards c
        JOIN auth.users u ON c.user_id = u.id
        WHERE u.username = $1
    `

	rows, err := tx.Query(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query user cards: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var card profile.CardInfo
		if err := rows.Scan(&card.CardNumber, &card.CardHolder, &card.ExpirationDate, &card.Cvv, &card.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan card info: %w", err)
		}
		cards = append(cards, card)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return cards, nil
}

// deleteCard removes a specific card associated with a user from the database.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - tx: The database transaction interface for executing the query.
// - username: The username associated with the card to be deleted.
// - cardNumber: The card number to be deleted.
//
// Returns:
// - error: An error if the operation fails or if no matching card is found.
//
// Workflow:
//  1. Executes a `DELETE` query to remove the card with the specified `card_number`
//     and `user_id` (determined from the `username`).
//  2. Checks the `RowsAffected()` count from the query result to ensure a card was deleted.
//     - If no rows were affected, returns a "card not found" error.
//  3. Returns any error encountered during query execution or validation.
//
// Error Handling:
// - Returns a "failed to delete card info" error if the query execution fails.
// - Returns a "card not found" error if no rows match the specified `username` and `cardNumber`.
func deleteCard(ctx context.Context, tx pgx.Tx, username, cardNumber string) error {
	const query = `
        DELETE FROM auth.cards
        WHERE user_id = (
            SELECT id FROM auth.users WHERE username = $1
        )
        AND card_number = $2
    `

	cmdTag, err := tx.Exec(ctx, query, username, cardNumber)
	if err != nil {
		return fmt.Errorf("failed to delete card info: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("card not found")
	}

	return nil
}
