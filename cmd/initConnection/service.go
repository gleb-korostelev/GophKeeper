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
func initServices(db db.IAdapter, key []byte) (
	profileSvc handler.ProfileSvc,
	authSvc handler.AuthSvc,
) {
	profileSvc = profile.NewService(db)
	authSvc = auth.NewService(db, key)

	return
}
