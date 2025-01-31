// Package profile defines the data structures related to user profile management,
// including card information and related metadata.
package profile

import "time"

// CardInfo represents the structure for storing information about a user's card.
type CardInfo struct {
	Username       string    `json:"username" validate:"required,min=3,max=50" example:"john_doe"`
	CardNumber     string    `json:"card_number" validate:"required,len=16,numeric" example:"1234567812345678"`
	CardHolder     string    `json:"card_holder" validate:"required,min=3,max=100" example:"John Doe"`
	ExpirationDate time.Time `json:"expiration_date" validate:"required" example:"2025-01-01"`
	Cvv            string    `json:"cvv" validate:"required,len=3,numeric" example:"123"`
	Metadata       string    `json:"metadata,omitempty" validate:"max=1000" example:"additional info"`
}
