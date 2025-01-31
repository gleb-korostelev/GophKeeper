package repository

import (
	"context"

	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetAccountByUserName(ctx context.Context, tx pgx.Tx, username string) (models.Account, error)
	UploadCardInfo(ctx context.Context, tx pgx.Tx, profile profile.CardInfo) error
	GetUserCards(ctx context.Context, tx pgx.Tx, username string) ([]profile.CardInfo, error)
	DeleteCard(ctx context.Context, tx pgx.Tx, username, cardNumber string) error
	InsertAccount(ctx context.Context, tx pgx.Tx, username string, secret []byte) (err error)
	UpdateAccountType(ctx context.Context, tx pgx.Tx, username string, accType models.AccountType) (err error)
}
