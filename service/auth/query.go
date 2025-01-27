package auth

import (
	"context"

	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/jackc/pgx/v5"
)

// getAccountByUserName retrieves an account record from the database by username.
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

// insertAccount inserts a new account record into the database.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - tx: The database transaction interface.
// - username: The username for the new account.
// - secret: The hashed password for the account.
//
// Returns:
// - err: An error if the insertion fails.
//
// Workflow:
// 1. Executes an `INSERT` query to add the account with the provided username and secret.
// 2. Returns any error encountered.

func insertAccount(ctx context.Context, tx pgx.Tx, username string, secret []byte) (err error) {
	const query = `
		INSERT INTO auth.users(username, secret)
		VALUES ($1, $2);
	`

	_, err = tx.Exec(ctx, query, username, secret)
	return
}

// updateAccountType updates the account type for a specific username.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - tx: The database transaction interface.
// - username: The username of the account to update.
// - accType: The new account type, represented as a `models.AccountType`.
//
// Returns:
// - err: An error if the update operation fails.
//
// Workflow:
// 1. Executes an `UPDATE` query to modify the `account_type`, `role_changed_at`, and `updated_at` fields.
// 2. Returns any error encountered.
func updateAccountType(ctx context.Context, tx pgx.Tx, username string, accType models.AccountType) (err error) {
	const query = `
		UPDATE auth.users
		SET account_type = $2,
		    role_changed_at = now(),
			updated_at = now()
		WHERE username = $1;
	`

	_, err = tx.Exec(ctx, query, username, accType)

	return
}
