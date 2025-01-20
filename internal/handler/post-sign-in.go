package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// PostSignInReq represents the structure of the request body for user sign-in.
//
// Fields:
// - Username: The username of the user attempting to sign in.
// - Password: The user's password.
// - Challenge: A challenge string used to validate the sign-in process.
type PostSignInReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Challenge string `json:"challenge"`
}

// PostSignInResp represents the structure of the response body for user sign-in.
//
// Fields:
// - Token: The generated authentication token for the user.
// - RefreshToken: A refresh token for renewing the authentication session.
type PostSignInResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// PostSignIn handles user authentication and generates access tokens.
// This function validates the provided credentials and challenge, then issues tokens for authenticated users.
//
// Workflow:
// 1. Parses and validates the request body to extract username, password, and challenge.
// 2. Verifies that the username, password, and challenge are not empty.
// 3. Calls the `AuthSvc.SignIn` service to authenticate the user and generate tokens.
// 4. Responds with the generated tokens in a structured JSON response.
//
// Parameters:
// - rw: The HTTP response writer used to send the response.
// - r: The HTTP request containing the request body and context.
//
// Response:
//   - HTTP Status Code: 200 OK on success.
//   - JSON Body:
//     {
//     "token": "<access_token>",
//     "refresh_token": "<refresh_token>"
//     }
//
// Error Handling:
// - 400 Bad Request: If the request body is invalid or malformed.
// - 400 Bad Request: If any required field (username, password, challenge) is missing or empty.
// - 401 Unauthorized: If authentication fails due to invalid credentials or challenge.
// - 500 Internal Server Error: If an unexpected error occurs.
//
// Example usage in a router setup:
//
//	router.HandleFunc("/api/v1/login", handler.PostSignIn).Methods("POST")
func (i *Implementation) PostSignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request body to extract username, password, and challenge.
	req, err := decoder.DecodeJson[PostSignInReq](r.Body)
	if err != nil {
		// Handle invalid JSON syntax or unexpected characters in the request body.
		if _, ok := err.(*json.SyntaxError); ok || strings.Contains(err.Error(), "invalid character") {
			handleErrResponse(rw, errInvalidRequestBody)
		} else {
			handleErrResponse(rw, err)
		}
		return
	}

	// Validate that all required fields are provided.
	if len(req.Username) == 0 || len(req.Password) == 0 || len(req.Challenge) == 0 {
		handleErrResponse(rw, errInvalidArgument)
		return
	}

	// Create a Profile object with the provided username and password.
	p := models.Profile{
		Username: req.Username,
		Password: req.Password,
	}

	// Authenticate the user and generate tokens using the authentication service.
	token, rToken, err := i.AuthSvc.SignIn(ctx, p, req.Challenge)
	if err != nil {
		handleErrResponse(rw, errAuthFailed)
		return
	}

	// Respond with the generated tokens.
	response.OK(rw, repackPostSignIn(token, rToken))
}

// repackPostSignIn converts the generated tokens into the API response format.
//
// Parameters:
// - token: The generated authentication token.
// - rToken: The generated refresh token.
//
// Returns:
// - A PostSignInResp structure containing the tokens.
//
// Example usage:
//
//	resp := repackPostSignIn("example_token", "example_refresh_token")
func repackPostSignIn(token, rToken string) PostSignInResp {
	return PostSignInResp{Token: token, RefreshToken: rToken}
}
