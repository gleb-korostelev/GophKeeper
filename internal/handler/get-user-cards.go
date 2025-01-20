package handler

import (
	"net/http"
	"time"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
	"github.com/gleb-korostelev/GophKeeper/middleware"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
)

// UsernameParam is the query parameter key for filtering by username in API endpoints.
const (
	UsernameParam = "username"
)

// GetUserCardsResp represents the structure of the API response for retrieving user cards.
//
// Fields:
// - Username: The username associated with the cards.
// - Cards: A slice of CardResp containing the user's card details.
type GetUserCardsResp struct {
	Username string     `json:"username"`
	Cards    []CardResp `json:"cards"`
}

// CardResp represents the structure of a single card's information in the response.
//
// Fields:
// - CardNumber: The card number (e.g., 16-digit card number).
// - CardHolder: The name of the cardholder.
// - ExpirationDate: The expiration date of the card.
// - Cvv: The CVV security code of the card.
// - Metadata: Additional metadata associated with the card.
type CardResp struct {
	CardNumber     string    `json:"card_number"`
	CardHolder     string    `json:"card_holder"`
	ExpirationDate time.Time `json:"expiration_date"`
	Cvv            string    `json:"cvv"`
	Metadata       string    `json:"metadata"`
}

// GetUserCards handles the retrieval of a user's saved card information.
// It performs the following steps:
// 1. Validates the user's token to identify the issuer.
// 2. Checks that the user has the required authorization level.
// 3. Retrieves the user's cards from the ProfileSvc service.
// 4. Returns the list of cards in a structured JSON response.
//
// Parameters:
// - rw: The HTTP response writer for sending the response.
// - r: The HTTP request containing context and headers.
//
// Workflow:
// - Validates the user's token using `middleware.GetIssuer`.
// - Retrieves the user's account using `AuthSvc.GetAccountByUserName`.
// - Ensures the user has an authorized account type.
// - Calls `ProfileSvc.GetUserCards` to fetch the cards.
// - Responds with a structured JSON response containing the cards.
//
// Error Handling:
// - 401 Unauthorized: If the token is invalid or missing.
// - 403 Forbidden: If the user lacks sufficient rights.
// - 500 Internal Server Error: If an error occurs while fetching cards.
//
// Example usage in a router setup:
//
//	router.HandleFunc("/api/v1/cards", handler.GetUserCards).Methods("GET")
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
//
// Parameters:
// - username: The username associated with the cards.
// - cards: A slice of CardInfo containing the raw card data.
//
// Returns:
// - A GetUserCardsResp structure containing the username and formatted card data.
//
// Example usage:
//
//	formattedResponse := repackGetCards("user123", []profile.CardInfo{...})
func repackGetCards(username string, cards []profile.CardInfo) GetUserCardsResp {
	// Prepare a new slice for the formatted card data.
	newcards := make([]CardResp, 0, len(cards))

	// Convert each CardInfo to the response format.
	for _, card := range cards {
		newcards = append(newcards, CardResp{
			CardNumber:     card.CardNumber,
			CardHolder:     card.CardHolder,
			ExpirationDate: card.ExpirationDate,
			Cvv:            card.Cvv,
			Metadata:       card.Metadata,
		})
	}

	// Return the formatted response.
	return GetUserCardsResp{Username: username, Cards: newcards}
}
