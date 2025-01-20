// Package otp provides utilities for validating one-time passwords (OTPs) and user credentials.
package otp

import (
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword verifies the provided OTP and password against stored values.
//
// Parameters:
// - otpCurr: The current OTP value.
// - otpPrev: The previous OTP value (used for fallback validation).
// - password: The plaintext password provided by the user.
// - msg: A message containing the OTP appended to a UUID prefix.
// - secret: The hashed password (stored securely).
//
// Returns:
// - bool: `true` if the provided OTP and password are valid; otherwise, `false`.
//
// Workflow:
// 1. Extracts the OTP portion from the message by removing the UUID prefix.
// 2. Validates whether the extracted OTP matches either the current or previous OTP.
// 3. If the OTP matches, verifies the plaintext password against the hashed password using bcrypt.
//
// Usage:
//
//	  if VerifyPassword(currOtp, prevOtp, userPassword, receivedMsg, storedSecret) {
//	      fmt.Println("Password and OTP are valid!")
//	  } else {
//	      fmt.Println("Invalid credentials.")
//	}
//
// Example:
//
//	currOtp := "123456"
//	prevOtp := "654321"
//	password := "securepassword"
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	msg := uuid.NewV4().String() + currOtp
//	isValid := VerifyPassword(currOtp, prevOtp, password, msg, hashedPassword)
//
// Behavior:
// - Returns `false` if the extracted OTP does not match either `otpCurr` or `otpPrev`.
// - Uses `bcrypt.CompareHashAndPassword` for secure password verification.
func VerifyPassword(otpCurr, otpPrev, password, msg string, secret []byte) bool {
	// Remove the UUID prefix from the message to extract the OTP.
	cutPrefix := len(uuid.Nil.String())

	// Check if the extracted OTP matches either the current or previous OTP.
	if msg[cutPrefix:] != otpCurr && msg[cutPrefix:] != otpPrev {
		return false
	}

	// Compare the provided password with the stored hashed password.
	err := bcrypt.CompareHashAndPassword(secret, []byte(password))
	return err == nil
}
