// Package auth provides services for managing user authentication, including profile creation,
// OTP-based challenge generation, and JWT-based sign-in.

package auth

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/pkg/claims"
	"github.com/gleb-korostelev/GophKeeper/pkg/otp"
	svc "github.com/gleb-korostelev/GophKeeper/service"
	"github.com/gleb-korostelev/GophKeeper/tools/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	sevenDays = time.Hour * 24 * 7 // Token expiration duration for refresh tokens.
	fiveMin   = time.Minute * 5    // Token expiration duration for short-lived tokens.
)

// service defines the implementation of the authentication service.
//
// Fields:
// - privateKey: The private key used for signing JWT tokens.
// - db: The database adapter for executing database operations.
type service struct {
	privateKey ed25519.PrivateKey
	db         db.IAdapter
}

// NewService creates a new instance of the authentication service.
//
// Parameters:
// - db: The database adapter for database operations.
// - privateKey: The private key for signing JWT tokens.
//
// Returns:
// - *service: The authentication service implementation.
func NewService(db db.IAdapter, privateKey ed25519.PrivateKey) *service {
	return &service{db: db, privateKey: privateKey}
}

// CreateProfile creates a new user profile or retrieves an existing one, returning an OTP challenge.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - profile: The user profile containing the username and password.
//
// Returns:
// - challenge: A combined UUID prefix and OTP challenge for verification.
// - err: An error if the operation fails.
//
// Workflow:
// 1. Attempts to retrieve an existing account by username.
// 2. If no account exists, creates a new account with a hashed password.
// 3. Generates an OTP challenge using the account's secret.
func (s *service) CreateProfile(ctx context.Context, profile models.Profile) (challenge string, err error) {
	var acc models.Account
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		acc, err = getAccountByUserName(ctx, tx, profile.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = acc.GenerateSecret(profile.Password)
				if err != nil {
					return err
				}

				err = insertAccount(ctx, tx, profile.Username, acc.Secret)
				if err != nil {
					return fmt.Errorf("error in insertAccount: %w", err)
				}
			} else {
				return fmt.Errorf("error in getAccountByUserName: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	challengePrefix := uuid.New().String()
	challenge, _, err = otp.GetTotp(acc.Secret)
	return strings.Join([]string{challengePrefix, challenge}, ""), err
}

// GetChallenge generates an OTP challenge for an existing user profile.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - profile: The user profile containing the username.
//
// Returns:
// - challenge: A combined UUID prefix and OTP challenge for verification.
// - err: An error if the operation fails.
func (s *service) GetChallenge(ctx context.Context, profile models.Profile) (challenge string, err error) {
	var acc models.Account
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		acc, err = getAccountByUserName(ctx, tx, profile.Username)
		if err != nil {
			return fmt.Errorf("error in getAccountByUserName: %w", err)
		}
		return nil
	})
	if err != nil {
		return
	}

	challengePrefix := uuid.New().String()
	challenge, _, err = otp.GetTotp(acc.Secret)
	return strings.Join([]string{challengePrefix, challenge}, ""), err
}

// SignIn authenticates a user by validating their OTP and password, then returns JWT tokens.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - profile: The user profile containing the username and password.
// - challenge: The OTP challenge to validate.
//
// Returns:
// - token: A short-lived JWT token.
// - refresh: A long-lived refresh token.
// - err: An error if authentication fails.
func (s *service) SignIn(ctx context.Context, profile models.Profile, challenge string) (token, refresh string, err error) {
	var acc models.Account
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		acc, err = getAccountByUserName(ctx, tx, profile.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return svc.ErrAccountNotFound
			}
			return fmt.Errorf("error in getAccountByUserName: %w", err)
		}
		return nil
	})
	if err != nil {
		return
	}

	otpCurr, otpPrev, err := otp.GetTotp(acc.Secret)
	if err != nil {
		return "", "", err
	}

	if !otp.VerifyPassword(otpCurr, otpPrev, profile.Password, challenge, acc.Secret) {
		return "", "", svc.ErrIncorrectPassword
	}

	if acc.AccountType == models.AccountUnauthorizedUser {
		acc.AccountType = models.AccountAuthorizedUser
		err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
			err = updateAccountType(ctx, tx, acc.Username, acc.AccountType)
			if err != nil {
				return fmt.Errorf("error in updateAccountType: %w", err)
			}
			return nil
		})
		if err != nil {
			return
		}
	}

	roleFunc := getRole(acc.AccountType)
	abilities := []claims.Ability{roleFunc(profile.Username)}

	claims := claims.NewClaims(
		time.Hour,
		claims.Role{
			Name:      profile.Username,
			Global:    true,
			Abilities: claims.ToAbilities(abilities...),
		},
	)
	token, refresh, err = claims.Sign(s.privateKey, profile.Username, fiveMin)
	return
}

// GetAccountByUserName retrieves an account by username.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellations.
// - username: The username of the account to retrieve.
//
// Returns:
// - acc: The account details.
// - err: An error if the operation fails or the account is not found.
func (s *service) GetAccountByUserName(ctx context.Context, username string) (acc models.Account, err error) {
	err = s.db.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		acc, err = getAccountByUserName(ctx, tx, username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return svc.ErrAccountNotFound
			} else {
				return fmt.Errorf("error in getAccountByUserName: %w", err)
			}
		}
		return nil
	})
	return acc, err
}
