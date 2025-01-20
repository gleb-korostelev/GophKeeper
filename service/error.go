// Package service provides common error definitions used across the GophKeeper application.
//
// These errors are used to standardize the error handling for account and authorization-related
// operations, ensuring consistent behavior and messaging throughout the application.

package service

import "errors"

var (
	// ErrAccountNotFound indicates that an account with the specified details could not be found.
	ErrAccountNotFound = errors.New("account not found")

	// ErrIncorrectPassword indicates that the provided password does not match the stored account secret.
	ErrIncorrectPassword = errors.New("incorrect password")

	// ErrNotAuthorized indicates that the user does not have sufficient permissions for the requested operation.
	ErrNotAuthorized = errors.New("not authorized")
)
