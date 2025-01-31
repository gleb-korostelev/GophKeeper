package repository

import (
	"context"

	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/jackc/pgx/v5"
)

// InsertAccount inserts a new account record into the database.
func InsertAccount(ctx context.Context, tx pgx.Tx, username string, secret []byte) (err error) {
	const query = `
		INSERT INTO auth.users(username, secret)
		VALUES ($1, $2);
	`

	_, err = tx.Exec(ctx, query, username, secret)
	return
}

// UpdateAccountType updates the account type for a specific username.
func UpdateAccountType(ctx context.Context, tx pgx.Tx, username string, accType models.AccountType) (err error) {
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
