// Package handler provides utilities for handling HTTP requests and generating appropriate responses.
// It also includes error handling logic to streamline response generation for common errors.
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
//
// Parameters:
// - rw: The HTTP response writer used to send the response.
// - err: The error that occurred, which determines the HTTP status code and message.
//
// Workflow:
// - Logs the error using the logger utility.
// - Matches the error against predefined errors and sends a response with the corresponding HTTP status code.
// - If the error is not recognized, it defaults to a 500 Internal Server Error.
//
// Error Mapping:
// - errInvalidRequestBody -> 400 Bad Request
// - errHashingPassword -> 500 Internal Server Error
// - middleware.ErrTokenInvalid -> 401 Unauthorized
// - middleware.ErrNotEnoughRights -> 403 Forbidden
// - errAuthFailed -> 401 Unauthorized
// - Any other error -> 500 Internal Server Error
//
// Example usage in a handler:
//
//	func ExampleHandler(rw http.ResponseWriter, r *http.Request) {
//	    err := someFunction()
//	    if err != nil {
//	        handleErrResponse(rw, err)
//	        return
//	    }
//	}
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
