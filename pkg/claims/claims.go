// Package claims provides utilities for managing JWT tokens,
// including roles, abilities, and authentication claims for the GophKeeper application.

package claims

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Role defines a user's role in the system, including its name, global status, and associated abilities.
//
// Fields:
// - Name: The name of the role (e.g., "admin").
// - Global: Indicates whether the role applies globally across all resources.
// - Abilities: A map of abilities assigned to the role, organized by scope.
type Role struct {
	Name      string    `json:"name"`
	Global    bool      `json:"global"`
	Abilities Abilities `json:"abilities,omitempty"`
}

// Abilities is a map of ability names to their corresponding scopes.
type Abilities map[string][]string

// Ability represents a specific capability within a defined scope.
//
// Fields:
// - Name: The name of the ability (e.g., "edit").
// - Scope: The specific scope for the ability (e.g., "project1").
type Ability struct {
	Name  string `json:"name"`
	Scope string `json:"scope,omitempty"`
}

// Claims extends JWT's StandardClaims with a Role field to represent user claims.
//
// Fields:
// - StandardClaims: The standard JWT claims such as `IssuedAt`, `ExpiresAt`, `Issuer`, and `Subject`.
// - Role: The user's role and associated abilities.
type Claims struct {
	jwt.StandardClaims
	Role
}

// NewClaims creates a new `Claims` object with a specified duration and role.
//
// Parameters:
// - duration: The validity duration of the token.
// - role: The role to be assigned to the claims.
//
// Returns:
// - *Claims: The created claims object.
func NewClaims(duration time.Duration, role Role) *Claims {
	iat := time.Now()
	eat := iat.Add(duration)
	return &Claims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  iat.Unix(),
			ExpiresAt: eat.Unix(),
			Issuer:    "gophkeeper",
			Subject:   "gophkeeper-auth",
		},
	}
}

// Parse validates and parses a JWT token using the provided public key.
//
// Parameters:
// - token: The JWT token string to be parsed.
// - publicKey: The Ed25519 public key for signature verification.
//
// Returns:
// - error: An error if token validation or parsing fails.
func (claims *Claims) Parse(token string, publicKey ed25519.PublicKey) error {
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return err
	}

	if !t.Valid {
		return errors.New("unreachable")
	}
	return nil
}

// Sign generates a JWT token and refresh token using the provided private key.
//
// Parameters:
// - key: The Ed25519 private key for signing the token.
// - sub: The subject (e.g., username) to include in the refresh token.
// - rTokenExp: The expiration duration for the refresh token.
//
// Returns:
// - string: The signed JWT token.
// - string: The signed refresh token.
// - error: Any error encountered during token generation.
func (claims *Claims) Sign(key ed25519.PrivateKey, sub string, rTokenExp time.Duration) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	jwtToken, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.New(jwt.SigningMethodEdDSA)
	rtClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("token is not valid JWT format")
	}
	rtClaims["sub"] = sub
	rtClaims["exp"] = time.Now().Add(rTokenExp).Unix()
	rt, err := refreshToken.SignedString(key)
	if err != nil {
		return "", "", err
	}

	return jwtToken, rt, nil
}

// ToAbilities converts a list of `Ability` objects into an `Abilities` map.
//
// Parameters:
// - items: A variadic list of `Ability` objects.
//
// Returns:
// - Abilities: A map of ability names to their scopes.
func ToAbilities(items ...Ability) Abilities {
	result := make(Abilities)
	for _, ability := range items {
		scopes, ok := result[ability.Name]
		if !ok {
			scopes = make([]string, 0, 1)
		}
		if len(ability.Scope) > 0 {
			scopes = append(scopes, ability.Scope)
		}
		result[ability.Name] = scopes
	}
	return result
}

// Predefined Role Generators

// SuperAdminRole creates a superadmin role with a specific scope.
func SuperAdminRole(username string) Ability {
	return Ability{
		Name:  "superadmin",
		Scope: username,
	}
}

// AdminRole creates an admin role with a specific scope.
func AdminRole(username string) Ability {
	return Ability{
		Name:  "admin",
		Scope: username,
	}
}

// RoleAuthorized creates an authorized user role with a specific scope.
func RoleAuthorized(username string) Ability {
	return Ability{
		Name:  "user",
		Scope: username,
	}
}

// Includes checks whether the claims include the specified abilities.
//
// Parameters:
// - abilities: A list of `Ability` objects to check.
//
// Returns:
// - bool: True if all specified abilities are included, false otherwise.
func (claims *Claims) Includes(abilities ...Ability) bool {
	for _, role := range abilities {
		grantedScopes := set(claims.Abilities[role.Name])
		if _, ok := grantedScopes[role.Scope]; !ok {
			return false
		}
	}

	return true
}

// set converts a slice of values into a map for quick lookups.
//
// Parameters:
// - values: A slice of comparable values.
//
// Returns:
// - map[T]struct{}: A map with the values as keys and empty structs as values.
func set[T comparable](values []T) map[T]struct{} {
	res := map[T]struct{}{}
	if values == nil {
		return res
	}
	for _, v := range values {
		res[v] = struct{}{}
	}
	return res
}
