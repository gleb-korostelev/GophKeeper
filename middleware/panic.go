package middleware

import (
	"net/http"

	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	"go.uber.org/zap"
)

// PanicMid is a middleware that recovers from panics in HTTP handlers,
// logs the error, and returns a 500 Internal Server Error response to the client.
func PanicMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Recover from any panic and handle the error.
			err := recover()
			if err != nil {
				// Log the error along with the request URL and method.
				logger.Error(err.(error).Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))

				// Return an HTTP 500 Internal Server Error response to the client.
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
