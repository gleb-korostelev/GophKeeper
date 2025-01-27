package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// PostChallenge handles the generation of an authentication challenge for a user.
// This challenge is used in authentication flows to validate user identity.
func (i *Implementation) PostChallenge(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request body to extract the username.
	req, err := decoder.DecodeJson[models.PostChallengeReq](r.Body)
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
func repackPostChallenge(challenge string) models.PostChallengeResp {
	return models.PostChallengeResp{Challenge: challenge}
}
