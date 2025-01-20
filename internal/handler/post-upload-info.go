package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/middleware"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/gleb-korostelev/GophKeeper/tools/decoder"
)

// PostUploadInfoReq represents the structure of the request body for uploading card information.
//
// Fields:
// - CardNumber: The card number to upload.
// - CardHolder: The name of the cardholder.
// - ExpirationDate: The expiration date of the card.
// - Cvv: The CVV security code of the card.
// - Metadata: Optional metadata associated with the card.
type PostUploadInfoReq struct {
	CardNumber     string    `json:"card_number"`
	CardHolder     string    `json:"card_holder"`
	ExpirationDate time.Time `json:"expiration_date"`
	Cvv            string    `json:"cvv"`
	Metadata       string    `json:"metadata"`
}

// PostUploadInfo handles the uploading or updating of card information for an authenticated user.
//
// Workflow:
// 1. Validates the user's token using `middleware.GetIssuer`.
// 2. Ensures the user has the necessary authorization level.
// 3. Parses and validates the request body to extract card information.
// 4. Calls the `ProfileSvc.UploadInfo` service to upload or update the card information.
// 5. Responds with a success message upon successful upload.
//
// Parameters:
// - rw: The HTTP response writer used to send the response.
// - r: The HTTP request containing the request body and context.
//
// Response:
// - HTTP Status Code: 200 OK on success.
// - JSON Body: `null`
//
// Error Handling:
// - 401 Unauthorized: If the token is invalid or missing.
// - 403 Forbidden: If the user lacks sufficient rights.
// - 400 Bad Request: If the request body is invalid or malformed.
// - 500 Internal Server Error: If an error occurs during the upload.
//
// Example usage in a router setup:
//
//	router.HandleFunc("/api/v1/cards", handler.PostUploadInfo).Methods("POST")
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
	req, err := decoder.DecodeJson[PostUploadInfoReq](r.Body)
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
