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

// DeleteCardInfoReq represents the structure of the request body for deleting a card.
//
// Fields:
// - CardNumber: The card number to be deleted.
type DeleteCardInfoReq struct {
	CardNumber string `json:"card_number"`
}

// DeleteCardInfo handles the deletion of a user's card information.
// It performs the following steps:
// 1. Verifies the token from the request context to identify the user.
// 2. Ensures the user has sufficient rights to delete card information.
// 3. Parses and validates the request body to extract the card number.
// 4. Calls the ProfileSvc service to delete the card.
//
// Parameters:
// - rw: The HTTP response writer.
// - r: The HTTP request containing the context, headers, and body.
//
// Workflow:
// - Validates the user's token using `middleware.GetIssuer`.
// - Ensures the user has an authorized account type using `AuthSvc.GetAccountByUserName`.
// - Decodes the request body to extract the `card_number`.
// - Calls the `ProfileSvc.DeleteCard` method to delete the card.
// - Responds with a 200 OK status on success or an appropriate error status on failure.
//
// Error Handling:
// - 401 Unauthorized: If the token is invalid or missing.
// - 403 Forbidden: If the user lacks sufficient rights.
// - 400 Bad Request: If the request body is invalid or malformed.
// - 500 Internal Server Error: If the deletion fails for unexpected reasons.
//
// Example usage in a router setup:
//
//	router.HandleFunc("/api/v1/cards", handler.DeleteCardInfo).Methods("DELETE")
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
	req, err := decoder.DecodeJson[DeleteCardInfoReq](r.Body)
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
