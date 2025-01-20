package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// PostChallengeReq represents the structure of the request body for the PostChallenge endpoint.
//
// Fields:
// - Username: The username for which the challenge is being requested.
type PostChallengeReq struct {
	Username string `json:"username"`
}

// PostChallengeResp represents the structure of the response body for the PostChallenge endpoint.
//
// Fields:
// - Challenge: A string containing the generated authentication challenge.
type PostChallengeResp struct {
	Challenge string `json:"challenge"`
}

// PostChallenge handles the generation of an authentication challenge for a user.
// This challenge is used in authentication flows to validate user identity.
//
// Workflow:
// 1. Parses the request body to extract the username.
// 2. Creates a `models.Profile` object with the provided username.
// 3. Calls the `AuthSvc.GetChallenge` service to generate a challenge for the user.
// 4. Responds with the generated challenge in a JSON response.
//
// Parameters:
// - rw: The HTTP response writer for sending the response.
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
// - 500 Internal Server Error: If the challenge generation fails.
//
// Example usage in a router setup:
//
//	router.HandleFunc("/api/v1/challenge", handler.PostChallenge).Methods("POST")
func (i *Implementation) PostChallenge(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request body to extract the username.
	req, err := decoder.DecodeJson[PostChallengeReq](r.Body)
	if err != nil {
		// Handle invalid JSON syntax or unexpected characters in the request body.
		if _, ok := err.(*json.SyntaxError); ok || strings.Contains(err.Error(), "invalid character") {
			handleErrResponse(rw, errInvalidRequestBody)
		} else {
			handleErrResponse(rw, err)
		}
		return
	}

	// Create a Profile object with the provided username.
	p := models.Profile{
		Username: req.Username,
		Password: "",
	}

	// Generate the challenge using the authentication service.
	challenge, err := i.AuthSvc.GetChallenge(ctx, p)
	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	// Respond with the generated challenge.
	response.OK(rw, repackPostChallenge(challenge))
}

// repackPostChallenge converts the generated challenge into the API response format.
//
// Parameters:
// - challenge: The string containing the generated authentication challenge.
//
// Returns:
// - A PostChallengeResp structure containing the challenge.
//
// Example usage:
//
//	resp := repackPostChallenge("example_challenge")
func repackPostChallenge(challenge string) PostChallengeResp {
	return PostChallengeResp{Challenge: challenge}
}
