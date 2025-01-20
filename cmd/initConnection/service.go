// Package initConnection provides functions to initialize core application components,
// including HTTP handlers, services, and middleware for GophKeeper.
package initConnection

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/config"
	"github.com/gleb-korostelev/GophKeeper/internal/handler"
	"github.com/gleb-korostelev/GophKeeper/internal/router"
	"github.com/gleb-korostelev/GophKeeper/service/auth"
	"github.com/gleb-korostelev/GophKeeper/service/profile"
	"github.com/gleb-korostelev/GophKeeper/tools/db"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	"github.com/rs/cors"
)

// InitImpl initializes the main HTTP handler for the GophKeeper application.
//
// It configures and initializes the following components:
// - Profile and Authentication services.
// - HTTP API handler with routing and middleware.
// - CORS middleware for cross-origin requests.
//
// The function uses configuration values from the application's settings:
// - config.IsSwaggerCreated: whether Swagger documentation should be included in the API.
// - config.JwtKey: the private key for signing JWT tokens (hex-encoded).
//
// Parameters:
// - ctx: a context.Context to manage the lifecycle of services.
// - adapter: a database adapter implementing the db.IAdapter interface.
// - port: the port number for the HTTP server.
//
// Returns:
// - http.Handler: the fully initialized HTTP handler, ready for use.
//
// Panics:
// - If the private key for signing JWT tokens is invalid or cannot be decoded.
//
// Example usage:
//
//	ctx := context.Background()
//	adapter := initConnection.NewDBConn(ctx)
//	httpHandler := initConnection.InitImpl(ctx, adapter, 8080)
//	http.ListenAndServe(":8080", httpHandler)
func InitImpl(
	ctx context.Context,
	adapter db.IAdapter,
	port int,
) http.Handler {

	isSwaggerCreated := config.GetConfigBool(config.IsSwaggerCreated)

	keyRaw := config.GetConfigString(config.JwtKey)
	keyBytes, err := hex.DecodeString(keyRaw)
	if err != nil {
		logger.Fatalf("error in hex.DecodeString: %v", err)
	}
	if l := len(keyBytes); l != ed25519.PrivateKeySize {
		logger.Fatalf("ed25519: bad private key length: %d", l)
	}

	profileSvc, authSvc := initServices(adapter, keyBytes)

	api := handler.NewImplementation(profileSvc, authSvc)
	r := router.CreateRouter(api, port, keyBytes, isSwaggerCreated)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodDelete,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions,
		},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Content-Length", "Content-Type", "Access-Control-Allow-Origin"},
	})

	return c.Handler(r)
}

// initServices initializes and returns the Profile and Authentication services.
//
// Parameters:
// - db: a database adapter implementing the db.IAdapter interface.
// - key: a byte slice representing the private key for signing JWT tokens.
//
// Returns:
// - profileSvc: an implementation of the handler.ProfileSvc interface.
// - authSvc: an implementation of the handler.AuthSvc interface.
//
// Example usage:
//
//	profileSvc, authSvc := initServices(adapter, key)
func initServices(db db.IAdapter, key []byte) (
	profileSvc handler.ProfileSvc,
	authSvc handler.AuthSvc,
) {
	profileSvc = profile.NewService(db)
	authSvc = auth.NewService(db, key)

	return
}
