package models

import "time"

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

// PostChallengeResp represents the structure of the response body for the PostChallenge endpoint.
//
// Fields:
// - Challenge: A string containing the generated authentication challenge.
type PostChallengeResp struct {
	Challenge string `json:"challenge"`
}

// PostProfileResp represents the structure of the response body for creating a user profile.
//
// Fields:
// - Challenge: A string containing the authentication challenge for the new profile.
type PostProfileResp struct {
	Challenge string `json:"challenge"`
}

// PostSignInResp represents the structure of the response body for user sign-in.
//
// Fields:
// - Token: The generated authentication token for the user.
// - RefreshToken: A refresh token for renewing the authentication session.
type PostSignInResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
