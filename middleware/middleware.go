// Package middleware provides tools for managing request authentication and authorization,
// including token validation and role-based access control.

package middleware

import (
	"context"
	"crypto/ed25519"
	"errors"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	auth "github.com/gleb-korostelev/GophKeeper/pkg/claims"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	"go.uber.org/zap"
)

// ctxKey represents a custom type for context keys to avoid collisions with other packages.
type ctxKey uint8

// Context keys for storing user-specific information.
const (
	CtxKeyUserID ctxKey = iota // The key for storing the user's ID.
	ctxKeyRoles                // The key for storing the user's roles or abilities.
)

// CoreMW represents the core middleware for handling authentication and authorization.
//
// Fields:
// - allowFake: A boolean to enable or disable fake authentication (for development or testing).
// - publicKey: The public key used for verifying JWT tokens.
type CoreMW struct {
	allowFake bool
	publicKey *ed25519.PublicKey
}

// NewCoreMW creates a new instance of CoreMW.
func NewCoreMW(allowFake bool, publicKey *ed25519.PublicKey) *CoreMW {
	return &CoreMW{
		allowFake: allowFake,
		publicKey: publicKey,
	}
}

// Predefined user roles.
const (
	RoleAdmin         = "admin"
	RoleSuperAdmin    = "superadmin"
	GuestUnauthorized = "unauthorized user"
	RoleAuthorized    = "authorized user"
)

// Header constants for request headers.
const (
	HeaderAuth = "Authorization" // Header for the authorization token.
)

// Errors related to authentication and authorization.
var (
	ErrTokenInvalid    = errors.New("bearer token is not correct")
	ErrNotEnoughRights = errors.New("not enough rights")
)

// Auth wraps an HTTP handler with authentication middleware.
func (a *CoreMW) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.allowFake {
			ctx, ok := a.FakeAuth(r)
			if ok {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		a.contextUpdate(next).ServeHTTP(w, r)
	}
}

// contextUpdate validates the token and updates the request context with user information.
func (a *CoreMW) contextUpdate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract the token from the Authorization header.
		bearerToken := r.Header.Get(HeaderAuth)
		tokenParts := strings.Split(bearerToken, " ")
		if len(tokenParts) != 2 {
			logger.Error("incorrect token")
			response.Unauthenticated(w, ErrTokenInvalid.Error())
			return
		}

		token := tokenParts[1]

		// Parse and validate the token.
		c, err := parseClaims(token, a.publicKey)
		if err != nil {
			logger.Error("error on contextUpdate.parseClaims not ok", zap.Error(err))
			response.Unauthenticated(w, ErrTokenInvalid.Error())
			return
		}

		// Update the context with user roles and address.
		ctx = context.WithValue(ctx, ctxKeyRoles, c.Role.Abilities)
		ctx = context.WithValue(ctx, CtxKeyUserID, c.Name)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FakeAuth simulates authentication for testing or development purposes.
func (a *CoreMW) FakeAuth(r *http.Request) (context.Context, bool) {
	ctx := r.Context()

	fakeAcc := r.Header.Get(HeaderAuth)
	fakeRoles := r.Header.Get(HeaderAuth)

	if len(fakeAcc) == 0 {
		return ctx, false
	}

	roles := strings.Split(fakeRoles, ",")

	ctx = context.WithValue(ctx, CtxKeyUserID, fakeAcc)
	ctx = context.WithValue(ctx, ctxKeyRoles, roles)
	return ctx, true
}

// parseClaims validates and parses a JWT token using the provided public key.
func parseClaims(token string, public *ed25519.PublicKey) (*auth.Claims, error) {
	claims := auth.Claims{}
	err := claims.Parse(token, *public)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

// GetIssuer retrieves the issuer (user identifier) from the context.
func GetIssuer(ctx context.Context) (string, error) {
	issuer, ok := ctx.Value(CtxKeyUserID).(string)
	if !ok || issuer == "" {
		return "", ErrTokenInvalid
	}
	return issuer, nil
}
