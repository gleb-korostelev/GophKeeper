// Package otp provides utilities for validating one-time passwords (OTPs) and user credentials.
package otp

import (
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword verifies the provided OTP and password against stored values.
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
