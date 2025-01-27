package database

import (
	"context"

	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/jackc/pgx/v5"
)

// GetAccountByUserName retrieves account information by username.
func GetAccountByUserName(ctx context.Context, tx pgx.Tx, username string) (acc models.Account, err error) {
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
