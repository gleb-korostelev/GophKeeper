// Package models defines the core data structures and utilities used across the application,
// including user accounts and profiles.

package models

import (
	"time"

	"github.com/gleb-korostelev/GophKeeper/middleware"
	"golang.org/x/crypto/bcrypt"
)

// AccountType defines the different types of user accounts.
//
// Constants:
// - AccountUnauthorizedUser: Represents an unauthorized user account.
// - AccountAuthorizedUser: Represents an authorized user account.
// - AccountRoleAdmin: Represents an admin user account.
// - AccountRoleSuperAdmin: Represents a superadmin user account.
type AccountType uint8

const (
	AccountUnauthorizedUser AccountType = iota + 1
	AccountAuthorizedUser
	AccountRoleAdmin
	AccountRoleSuperAdmin
)

// rolesHumanReadable maps account types to human-readable role names.
// These role names are defined in the `middleware` package.
var rolesHumanReadable = map[AccountType]string{
	AccountUnauthorizedUser: middleware.RoleUnauthorized,
	AccountAuthorizedUser:   middleware.RoleAuthorized,
	AccountRoleAdmin:        middleware.RoleAdmin,
	AccountRoleSuperAdmin:   middleware.RoleSuperAdmin,
}

// Account represents a user account in the system.
//
// Fields:
// - ID: The unique identifier for the account.
// - Username: The username associated with the account.
// - Secret: The hashed password for the account, stored as a byte slice.
// - AccountType: The type of the account, represented as an `AccountType`.
// - CreatedAt: The timestamp when the account was created.
// - RoleChangedAt: The timestamp when the account's role was last changed.
// - UpdatedAt: The timestamp when the account was last updated.
type Account struct {
	ID            int
	Username      string
	Secret        []byte
	AccountType   AccountType
	CreatedAt     time.Time
	RoleChangedAt time.Time
	UpdatedAt     time.Time
}

// Profile represents the basic profile information for user authentication.
//
// Fields:
// - Username: The username of the user.
// - Password: The plaintext password of the user.
type Profile struct {
	Username string
	Password string
}

// GenerateSecret hashes a plaintext password and assigns it to the account's Secret field.
//
// Parameters:
// - password: The plaintext password to hash.
//
// Returns:
// - error: An error if password hashing fails, otherwise nil.
//
// Workflow:
// 1. Uses bcrypt to hash the provided password with a default cost.
// 2. Assigns the hashed password to the `Secret` field of the `Account`.
//
// Example usage:
//
//	account := Account{}
//	err := account.GenerateSecret("securepassword")
//	if err != nil {
//	    log.Fatalf("Failed to generate secret: %v", err)
//	}
func (a *Account) GenerateSecret(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a.Secret = hashedPassword
	return nil
}
