package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// PostCreateProfileReq represents the structure of the request body for creating a user profile.
//
// Fields:
// - Username: The desired username for the new profile.
// - Password: The password associated with the profile.
type PostCreateProfileReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostProfileResp represents the structure of the response body for creating a user profile.
//
// Fields:
// - Challenge: A string containing the authentication challenge for the new profile.
type PostProfileResp struct {
	Challenge string `json:"challenge"`
}

// PostCreateProfile handles the creation of a new user profile and generates an authentication challenge.
//
// Workflow:
// 1. Parses the request body to extract the username and password.
// 2. Validates the parsed data.
// 3. Calls the `AuthSvc.CreateProfile` service to create the user profile and generate a challenge.
// 4. Responds with the generated challenge in a structured JSON response.
//
// Parameters:
// - rw: The HTTP response writer used to send the response.
// - r: The HTTP request containing the request body and context.
//
// Response:
//   - HTTP Status Code: 200 OK on success.
//   - JSON Body:
//     {
//     "challenge": "<generated_challenge>"
//     }
//
// Error Handling:
// - 400 Bad Request: If the request body is invalid or malformed.
// - 500 Internal Server Error: If profile creation or challenge generation fails.
//
// Example usage in a router setup:
//
//	router.HandleFunc("/api/v1/register", handler.PostCreateProfile).Methods("POST")
func (i *Implementation) PostCreateProfile(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request body to extract the username and password.
	req, err := decoder.DecodeJson[PostCreateProfileReq](r.Body)
	if err != nil {
		// Handle invalid JSON syntax or unexpected characters in the request body.
		if _, ok := err.(*json.SyntaxError); ok || strings.Contains(err.Error(), "invalid character") {
			handleErrResponse(rw, errInvalidRequestBody)
		} else {
			handleErrResponse(rw, err)
		}
		return
	}

	// Create a Profile object with the provided username and password.
	p := models.Profile{
		Username: req.Username,
		Password: req.Password,
	}

	// Create the profile and generate the challenge using the authentication service.
	challenge, err := i.AuthSvc.CreateProfile(ctx, p)
	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	// Respond with the generated challenge.
	response.OK(rw, repackPostCreateProfile(challenge))
}

// repackPostCreateProfile converts the generated challenge into the API response format.
//
// Parameters:
// - challenge: The string containing the generated authentication challenge.
//
// Returns:
// - A PostProfileResp structure containing the challenge.
//
// Example usage:
//
//	resp := repackPostCreateProfile("example_challenge")
func repackPostCreateProfile(challenge string) PostProfileResp {
	return PostProfileResp{Challenge: challenge}
}
