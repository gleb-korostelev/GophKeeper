package profile

import (
	"context"
	"fmt"

	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/jackc/pgx/v5"
)

// getAccountByUserName retrieves account information by username.
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
