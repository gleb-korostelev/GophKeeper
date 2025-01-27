package handler

import (
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/internal/handler/response"
)

// Healthcheck provides a simple endpoint to verify the health status of the application.
// It is commonly used for monitoring and troubleshooting purposes.
func (i *Implementation) Healthcheck(rw http.ResponseWriter, r *http.Request) {
	response.Healthcheck(rw)
}
