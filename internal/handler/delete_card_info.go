// Package handler provides HTTP handlers for managing user actions, such as deleting card information.
package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/middleware"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// DeleteCardInfo handles the deletion of a user's card information.
func (i *Implementation) DeleteCardInfo(rw http.ResponseWriter, r *http.Request) {
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

	// Decode the request body to extract the card number.
	req, err := decoder.DecodeJson[models.DeleteCardInfoReq](r.Body)
	if err != nil {
		if _, ok := err.(*json.SyntaxError); ok || strings.Contains(err.Error(), "invalid character") {
			handleErrResponse(rw, errInvalidRequestBody)
		} else {
			handleErrResponse(rw, err)
		}
		return
	}

	// Call the Profile service to delete the card for the user.
	err = i.ProfileSvc.DeleteCard(ctx, acc.Username, req.CardNumber)
	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	// Respond with a success message.
	response.OK(rw, nil)
}
