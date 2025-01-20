// Package profile defines the data structures related to user profile management,
// including card information and related metadata.
package profile

import "time"

// CardInfo represents the structure for storing information about a user's card.
//
// Fields:
// - Username: The username of the card owner.
//   - JSON: "username"
//   - Validation: Required, minimum length 3, maximum length 50.
//   - Example: "john_doe"
//
// - CardNumber: The card number.
//   - JSON: "card_number"
//   - Validation: Required, 16 numeric characters.
//   - Example: "1234567812345678"
//
// - CardHolder: The name of the cardholder as printed on the card.
//   - JSON: "card_holder"
//   - Validation: Required, minimum length 3, maximum length 100.
//   - Example: "John Doe"
//
// - ExpirationDate: The card's expiration date.
//   - JSON: "expiration_date"
//   - Validation: Required.
//   - Example: "2025-01-01"
//   - Type: time.Time
//
// - Cvv: The card's CVV security code.
//   - JSON: "cvv"
//   - Validation: Required, exactly 3 numeric characters.
//   - Example: "123"
//
// - Metadata: Optional metadata or additional information about the card.
//   - JSON: "metadata"
//   - Validation: Maximum length 1000 characters.
//   - Example: "additional info"
type CardInfo struct {
	Username       string    `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
	CardNumber     string    `json:"card_number" validate:"required,len=16,numeric" example:"1234567812345678"`
	CardHolder     string    `json:"card_holder" validate:"required,min=3,max=100" example:"John Doe"`
	ExpirationDate time.Time `json:"expiration_date" validate:"required" example:"2025-01-01"`
	Cvv            string    `json:"cvv" validate:"required,len=3,numeric" example:"123"`
	Metadata       string    `json:"metadata,omitempty" validate:"max=1000" example:"additional info"`
}
