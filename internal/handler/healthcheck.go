package handler

import (
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
)

// Healthcheck provides a simple endpoint to verify the health status of the application.
// It is commonly used for monitoring and troubleshooting purposes.
//
// Workflow:
// - Responds with a 200 OK HTTP status and a JSON body indicating the application is "healthy."
//
// Response Structure:
//   - HTTP Status Code: 200 OK
//   - JSON Body:
//     {
//     "status": "healthy"
//     }
//
// Example usage in a router setup:
//
//	router.HandleFunc("/healthcheck", handler.Healthcheck).Methods("GET")
//
// Parameters:
// - rw: The HTTP response writer for sending the response.
// - r: The HTTP request, though it is unused in this function.
func (i *Implementation) Healthcheck(rw http.ResponseWriter, r *http.Request) {
	response.Healthcheck(rw)
}
