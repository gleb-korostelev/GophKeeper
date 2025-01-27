// Package response provides utilities for structuring HTTP responses in a consistent
// JSON format. It simplifies setting HTTP status codes, messages, and response data.
package response

import (
	"encoding/json"
	"net/http"
)

// Response represents the structure of an HTTP response in JSON format.
type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

// HealthcheckResp represents the response structure for health check endpoints.
type HealthcheckResp struct {
	Status string `json:"status"`
}

// OK sends a 200 OK HTTP response with the provided data.
func OK(rw http.ResponseWriter, data interface{}) {
	result(rw, http.StatusOK, "Success", data)
}

// BadRequest sends a 400 Bad Request HTTP response with the provided error message.
func BadRequest(rw http.ResponseWriter, message string) {
	result(rw, http.StatusBadRequest, message, nil)
}

// Unauthenticated sends a 401 Unauthorized HTTP response with the provided error message.
func Unauthenticated(rw http.ResponseWriter, message string) {
	result(rw, http.StatusUnauthorized, message, nil)
}

// NotFound sends a 404 Not Found HTTP response with the provided error message.
func NotFound(rw http.ResponseWriter, message string) {
	result(rw, http.StatusNotFound, message, nil)
}

// Forbidden sends a 403 Forbidden HTTP response with the provided error message.
func Forbidden(rw http.ResponseWriter, message string) {
	result(rw, http.StatusForbidden, message, nil)
}

// Internal sends a 500 Internal Server Error HTTP response with the provided error message.
func Internal(rw http.ResponseWriter, message string) {
	result(rw, http.StatusInternalServerError, message, nil)
}

// result is a helper function to send an HTTP response in JSON format.
// It determines the success status based on the HTTP status code.
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
func Healthcheck(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(HealthcheckResp{
		Status: "healthy",
	})
}
