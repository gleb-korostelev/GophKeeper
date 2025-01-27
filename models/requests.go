package models

import "time"

// DeleteCardInfoReq represents the structure of the request body for deleting a card.
//
// Fields:
// - CardNumber: The card number to be deleted.
type DeleteCardInfoReq struct {
	CardNumber string `json:"card_number"`
}

// PostChallengeReq represents the structure of the request body for the PostChallenge endpoint.
//
// Fields:
// - Username: The username for which the challenge is being requested.
type PostChallengeReq struct {
	Username string `json:"username"`
}

// PostCreateProfileReq represents the structure of the request body for creating a user profile.
//
// Fields:
// - Username: The desired username for the new profile.
// - Password: The password associated with the profile.
type PostCreateProfileReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostSignInReq represents the structure of the request body for user sign-in.
//
// Fields:
// - Username: The username of the user attempting to sign in.
// - Password: The user's password.
// - Challenge: A challenge string used to validate the sign-in process.
type PostSignInReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Challenge string `json:"challenge"`
}

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
