package handler

import (
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/middleware"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
)

// UsernameParam is the query parameter key for filtering by username in API endpoints.
const (
	UsernameParam = "username"
)

// GetUserCards handles the retrieval of a user's saved card information.
func (i *Implementation) GetUserCards(rw http.ResponseWriter, r *http.Request) {
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

	// Retrieve the user's cards from the profile service.
	cards, err := i.ProfileSvc.GetUserCards(ctx, acc.Username)
	if err != nil {
		handleErrResponse(rw, err)
		return
	}

	// Send the response with the repacked card data.
	response.OK(rw, repackGetCards(acc.Username, cards))
}

// repackGetCards converts a slice of CardInfo to the API response structure (GetUserCardsResp).
func repackGetCards(username string, cards []profile.CardInfo) models.GetUserCardsResp {
	// Prepare a new slice for the formatted card data.
	newcards := make([]models.CardResp, 0, len(cards))

	// Convert each CardInfo to the response format.
	for _, card := range cards {
		newcards = append(newcards, models.CardResp{
			CardNumber:     card.CardNumber,
			CardHolder:     card.CardHolder,
			ExpirationDate: card.ExpirationDate,
			Cvv:            card.Cvv,
			Metadata:       card.Metadata,
		})
	}

	// Return the formatted response.
	return models.GetUserCardsResp{Username: username, Cards: newcards}
}
