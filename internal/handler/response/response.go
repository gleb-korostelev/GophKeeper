// Package response provides utilities for structuring HTTP responses in a consistent
// JSON format. It simplifies setting HTTP status codes, messages, and response data.
package response

import (
	"encoding/json"
	"net/http"
)

// Response represents the structure of an HTTP response in JSON format.
// It uses a generic type `T` to allow flexibility in the type of data being returned.
//
// Fields:
// - Success: A boolean indicating if the request was successful (true) or not (false).
// - Message: An optional string providing additional information about the response.
// - Data: A generic field containing the payload for the response, if any.
type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

// HealthcheckResp represents the response structure for health check endpoints.
//
// Fields:
// - Status: A string indicating the health status of the application.
type HealthcheckResp struct {
	Status string `json:"status"`
}

// OK sends a 200 OK HTTP response with the provided data.
//
// Parameters:
// - rw: The HTTP response writer.
// - data: The payload to include in the response.
//
// Example usage:
//
//	response.OK(rw, map[string]string{"key": "value"})
func OK(rw http.ResponseWriter, data interface{}) {
	result(rw, http.StatusOK, "Success", data)
}

// BadRequest sends a 400 Bad Request HTTP response with the provided error message.
//
// Parameters:
// - rw: The HTTP response writer.
// - message: A string describing the error.
//
// Example usage:
//
//	response.BadRequest(rw, "Invalid input data")
func BadRequest(rw http.ResponseWriter, message string) {
	result(rw, http.StatusBadRequest, message, nil)
}

// Unauthenticated sends a 401 Unauthorized HTTP response with the provided error message.
//
// Parameters:
// - rw: The HTTP response writer.
// - message: A string describing the error.
//
// Example usage:
//
//	response.Unauthenticated(rw, "Invalid credentials")
func Unauthenticated(rw http.ResponseWriter, message string) {
	result(rw, http.StatusUnauthorized, message, nil)
}

// NotFound sends a 404 Not Found HTTP response with the provided error message.
//
// Parameters:
// - rw: The HTTP response writer.
// - message: A string describing the error.
//
// Example usage:
//
//	response.NotFound(rw, "Resource not found")
func NotFound(rw http.ResponseWriter, message string) {
	result(rw, http.StatusNotFound, message, nil)
}

// Forbidden sends a 403 Forbidden HTTP response with the provided error message.
//
// Parameters:
// - rw: The HTTP response writer.
// - message: A string describing the error.
//
// Example usage:
//
//	response.Forbidden(rw, "Access denied")
func Forbidden(rw http.ResponseWriter, message string) {
	result(rw, http.StatusForbidden, message, nil)
}

// Internal sends a 500 Internal Server Error HTTP response with the provided error message.
//
// Parameters:
// - rw: The HTTP response writer.
// - message: A string describing the error.
//
// Example usage:
//
//	response.Internal(rw, "Internal server error")
func Internal(rw http.ResponseWriter, message string) {
	result(rw, http.StatusInternalServerError, message, nil)
}

// result is a helper function to send an HTTP response in JSON format.
// It determines the success status based on the HTTP status code.
//
// Parameters:
// - rw: The HTTP response writer.
// - status: The HTTP status code.
// - message: A string describing the response.
// - data: The payload to include in the response, if any.
//
// Example usage:
//
//	result(rw, http.StatusOK, "Success", map[string]string{"key": "value"})
func result(rw http.ResponseWriter, status int, message string, data any) {
	success := 200 <= status && status <= 300

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(Response[any]{
		Success: success,
		Message: message,
		Data:    data,
	})
}

// Healthcheck sends a 200 OK HTTP response indicating the health status of the application.
//
// Parameters:
// - rw: The HTTP response writer.
//
// Example usage:
//
//	response.Healthcheck(rw)
func Healthcheck(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(HealthcheckResp{
		Status: "healthy",
	})
}
