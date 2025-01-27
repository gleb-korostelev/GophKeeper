package handler

import (
	"errors"
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/middleware"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
)

// Predefined errors for common failure scenarios.
var (
	// errInvalidRequestBody indicates that the request body could not be parsed or is invalid.
	errInvalidRequestBody = errors.New("invalid request body")

	// errHashingPassword indicates a failure in the password hashing process.
	errHashingPassword = errors.New("error hashing password")

	// errInvalidArgument indicates that an argument provided in the request is invalid.
	errInvalidArgument = errors.New("invalid argument")

	// errAuthFailed indicates that user authentication has failed.
	errAuthFailed = errors.New("authentication failed")
)

// handleErrResponse sends an appropriate HTTP response based on the provided error.
// It maps predefined errors to specific HTTP status codes and response messages,
// ensuring consistent error handling across the application.
func handleErrResponse(rw http.ResponseWriter, err error) {
	defer logger.Info(err)

	switch err {
	case errInvalidRequestBody:
		// Handle invalid request body errors.
		response.BadRequest(rw, err.Error())
	case errHashingPassword:
		// Handle errors related to password hashing.
		response.Internal(rw, err.Error())
	case middleware.ErrTokenInvalid:
		// Handle token invalid errors (unauthorized access).
		response.Unauthenticated(rw, err.Error())
	case middleware.ErrNotEnoughRights:
		// Handle insufficient permission errors (forbidden access).
		response.Forbidden(rw, err.Error())
	case errAuthFailed:
		// Handle authentication failure errors (unauthorized access).
		response.Unauthenticated(rw, err.Error())
	default:
		// Default case for unrecognized errors.
		response.Internal(rw, err.Error())
	}
	return
}
