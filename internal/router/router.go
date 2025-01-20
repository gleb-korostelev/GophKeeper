package router

import (
	"crypto/ed25519"
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/internal/handler"
	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/middleware"
	"github.com/gleb-korostelev/GophKeeper/tools/swagger"
	"github.com/gorilla/mux"
)

// AppName is the name of the application, used in Swagger documentation.
const AppName = "gophkeeper"

// CreateRouter initializes the application's HTTP router, integrating handlers, middleware, and Swagger documentation.
//
// Parameters:
// - impl: An implementation of the `handler.API` interface containing the HTTP handlers for the application.
// - appPort: The main application port for the router.
// - authKey: The Ed25519 private key for authentication purposes.
// - isSwaggerCreated: A boolean indicating whether Swagger documentation has already been generated.
//
// Returns:
// - *mux.Router: A configured Gorilla Mux router with registered handlers, middleware, and Swagger support.
//
// Workflow:
// 1. Extracts the public key from the provided `authKey`.
// 2. Initializes the core middleware with the public key.
// 3. Registers handlers for all API endpoints with Swagger metadata.
// 4. Returns a new router initialized via `NewAPI`.
//
// Example usage:
//
//	router := CreateRouter(apiImplementation, 8080, authKey, false)
//
// Middleware:
// - The `mw.Auth` middleware is applied to endpoints requiring authentication.
// - Swagger headers are configured for authenticated endpoints.
//
// Swagger Integration:
// - Each handler is defined with Swagger metadata, including paths, request/response bodies, and options.
//
// Registered Endpoints:
// - `/healthcheck`: Checks the application's health.
// - `/api/v1/challenge`: Retrieves an authentication challenge.
// - `/api/v1/register`: Registers a new user profile.
// - `/api/v1/login`: Authenticates a user and issues tokens.
// - `/api/v1/upload-card-info`: Uploads or updates card information.
// - `/api/v1/cards` (GET): Retrieves all user cards.
// - `/api/v1/cards` (DELETE): Deletes a specific user card.
func CreateRouter(impl handler.API, appPort int, authKey ed25519.PrivateKey, isSwaggerCreated bool) *mux.Router {
	// Extract the public key for middleware initialization.
	pub := authKey.Public().(ed25519.PublicKey)

	// Initialize core middleware with the public key.
	mw := middleware.NewCoreMW(true, &pub)

	// Define handlers with Swagger metadata.
	var handlers = []swagger.Handler{
		{
			HandlerFunc:      impl.Healthcheck,
			Path:             "/healthcheck",
			Method:           http.MethodGet,
			Description:      "Healthcheck",
			ResponseBody:     response.HealthcheckResp{},
			ResponseMimeType: swagger.MimeJson,
			Opts:             []swagger.Option{},
			Tag:              Hc,
		},
		{
			HandlerFunc:  http.HandlerFunc(impl.PostChallenge),
			Path:         "/api/v1/challenge",
			Method:       http.MethodPost,
			Description:  "Get challenge for wallet",
			ResponseBody: response.Response[handler.PostChallengeResp]{},
			RequestBody:  handler.PostChallengeReq{},
		},
		{
			HandlerFunc:  http.HandlerFunc(impl.PostCreateProfile),
			Path:         "/api/v1/register",
			Method:       http.MethodPost,
			Description:  "Register system account",
			ResponseBody: response.Response[handler.PostProfileResp]{},
			RequestBody:  handler.PostCreateProfileReq{},
		},
		{
			HandlerFunc:  http.HandlerFunc(impl.PostSignIn),
			Path:         "/api/v1/login",
			Method:       http.MethodPost,
			Description:  "Login to accounts system",
			ResponseBody: response.Response[handler.PostSignInResp]{},
			RequestBody:  handler.PostSignInReq{},
		},
		{
			HandlerFunc:  mw.Auth(impl.PostUploadInfo),
			Path:         "/api/v1/upload-card-info",
			Method:       http.MethodPost,
			Description:  "Uploads or edits new card info",
			ResponseBody: response.Response[struct{}]{},
			RequestBody:  handler.PostUploadInfoReq{},
			Opts: []swagger.Option{
				swagger.HeaderOpt{
					Name:        middleware.HeaderAuth,
					Type:        swagger.String,
					Required:    true,
					Description: `Required 'Bearer ' prefix`,
				},
			},
		},
		{
			HandlerFunc:  mw.Auth(impl.GetUserCards),
			Path:         "/api/v1/cards",
			Method:       http.MethodGet,
			Description:  "Get card details",
			ResponseBody: response.Response[handler.GetUserCardsResp]{},
			Opts: []swagger.Option{
				swagger.HeaderOpt{
					Name:        middleware.HeaderAuth,
					Type:        swagger.String,
					Required:    true,
					Description: `Required 'Bearer ' prefix`,
				},
			},
		},
		{
			HandlerFunc:  mw.Auth(impl.DeleteCardInfo),
			Path:         "/api/v1/cards",
			Method:       http.MethodDelete,
			Description:  "Delete specific card",
			ResponseBody: response.Response[struct{}]{},
			RequestBody:  handler.DeleteCardInfoReq{},
			Opts: []swagger.Option{
				swagger.HeaderOpt{
					Name:        middleware.HeaderAuth,
					Type:        swagger.String,
					Required:    true,
					Description: `Required 'Bearer ' prefix`,
				},
			},
		},
	}

	// Create and return the new API router.
	return NewAPI(AppName, appPort, isSwaggerCreated, handlers)
}
