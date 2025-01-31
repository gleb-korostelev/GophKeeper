package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// PostCreateProfile handles the creation of a new user profile and generates an authentication challenge.
func (i *Implementation) PostCreateProfile(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request body to extract the username and password.
	req, err := decoder.DecodeJson[models.PostCreateProfileReq](r.Body)
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
func repackPostCreateProfile(challenge string) models.PostProfileResp {
	return models.PostProfileResp{Challenge: challenge}
}
