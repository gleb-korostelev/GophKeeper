package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// PostSignIn handles user authentication and generates access tokens.
// This function validates the provided credentials and challenge, then issues tokens for authenticated users.
func (i *Implementation) PostSignIn(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode the request body to extract username, password, and challenge.
	req, err := decoder.DecodeJson[models.PostSignInReq](r.Body)
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
func repackPostSignIn(token, rToken string) models.PostSignInResp {
	return models.PostSignInResp{Token: token, RefreshToken: rToken}
}
