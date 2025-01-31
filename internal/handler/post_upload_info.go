package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/middleware"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// PostUploadInfo handles the uploading or updating of card information for an authenticated user.
func (i *Implementation) PostUploadInfo(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Retrieve the issuer (user ID or token subject) from the request context.
	issuer, err := middleware.GetIssuer(ctx)
	if err != nil {
		handleErrResponse(rw, middleware.ErrTokenInvalid)
		return
	}

	// Retrieve the user's account details from the authentication service.
	var acc models.Account
	acc, err = i.AuthSvc.GetAccountByUserName(ctx, issuer)
	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	// Ensure the user has sufficient rights to perform this action.
	if acc.AccountType != models.AccountAuthorizedUser {
		handleErrResponse(rw, middleware.ErrNotEnoughRights)
		return
	}

	// Decode the request body to extract card information.
	req, err := decoder.DecodeJson[models.PostUploadInfoReq](r.Body)
	if err != nil {
		// Handle invalid JSON syntax or unexpected characters in the request body.
		if _, ok := err.(*json.SyntaxError); ok || strings.Contains(err.Error(), "invalid character") {
			handleErrResponse(rw, errInvalidRequestBody)
		} else {
			handleErrResponse(rw, err)
		}
		return
	}

	// Create a CardInfo object with the extracted data.
	p := profile.CardInfo{
		Username:       acc.Username,
		CardNumber:     req.CardNumber,
		CardHolder:     req.CardHolder,
		ExpirationDate: req.ExpirationDate,
		Cvv:            req.Cvv,
		Metadata:       req.Metadata,
	}

	// Upload the card information using the profile service.
	err = i.ProfileSvc.UploadInfo(ctx, p)
	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	// Respond with a success message.
	response.OK(rw, nil)
}
