package handler

import (
	"context"
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
)

// API defines the interface for the handler's API.
// Each method corresponds to an HTTP endpoint and performs a specific operation.
//
// Methods:
// - Healthcheck: Checks the health status of the application.
// - PostSignIn: Handles user sign-in and generates tokens.
// - PostCreateProfile: Creates a new user profile.
// - PostChallenge: Retrieves a challenge for user authentication.
// - PostUploadInfo: Uploads or updates card information for a user.
// - GetUserCards: Retrieves all cards associated with a user.
// - DeleteCardInfo: Deletes a specific card associated with a user.
type API interface {
	Healthcheck(rw http.ResponseWriter, r *http.Request)
	PostSignIn(rw http.ResponseWriter, r *http.Request)
	PostCreateProfile(rw http.ResponseWriter, r *http.Request)
	PostChallenge(rw http.ResponseWriter, r *http.Request)
	PostUploadInfo(rw http.ResponseWriter, r *http.Request)
	GetUserCards(rw http.ResponseWriter, r *http.Request)
	DeleteCardInfo(rw http.ResponseWriter, r *http.Request)
}

// ProfileSvc defines the interface for interacting with the profile service.
//
// Methods:
// - UploadInfo: Uploads or updates card information for a specific user.
// - GetUserCards: Retrieves all cards associated with a username.
// - DeleteCard: Deletes a specific card for a user based on username and card number.
type ProfileSvc interface {
	UploadInfo(ctx context.Context, profile profile.CardInfo) (err error)
	GetUserCards(ctx context.Context, username string) ([]profile.CardInfo, error)
	DeleteCard(ctx context.Context, username, cardNumber string) (err error)
}

// AuthSvc defines the interface for interacting with the authentication service.
//
// Methods:
// - CreateProfile: Creates a new user profile and generates a challenge for authentication.
// - GetChallenge: Retrieves an authentication challenge for a user.
// - SignIn: Authenticates a user and generates an access token and refresh token.
// - GetAccountByUserName: Retrieves account details for a specific username.
type AuthSvc interface {
	CreateProfile(ctx context.Context, profile models.Profile) (challenge string, err error)
	GetChallenge(ctx context.Context, profile models.Profile) (challenge string, err error)
	SignIn(ctx context.Context, profile models.Profile, challenge string) (token, refresh string, err error)
	GetAccountByUserName(ctx context.Context, username string) (acc models.Account, err error)
}

// Implementation provides the concrete implementation of the API interface.
// It acts as a bridge between the HTTP layer and the services.
//
// Fields:
// - ProfileSvc: The service responsible for managing user profiles.
// - AuthSvc: The service responsible for managing authentication.
type Implementation struct {
	ProfileSvc ProfileSvc
	AuthSvc    AuthSvc
}

// NewImplementation creates a new instance of the API implementation.
//
// Parameters:
// - profileSvc: The service for managing user profile operations.
// - authSvc: The service for managing authentication operations.
func NewImplementation(profileSvc ProfileSvc, authSvc AuthSvc) API {
	return &Implementation{
		ProfileSvc: profileSvc,
		AuthSvc:    authSvc,
	}
}
